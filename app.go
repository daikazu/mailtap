package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"mailtap/internal/models"
	"mailtap/internal/notify"
	"mailtap/internal/otp"
	"mailtap/internal/settings"
	"mailtap/internal/smtp"
	"mailtap/internal/storage"
)

type App struct {
	ctx        context.Context
	repo       *storage.Repository
	smtpServer *smtp.Server
	port       string
}

func NewApp(port string) *App {
	return &App{port: port}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dbPath := storage.DefaultDBPath()
	db, err := storage.NewDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	a.repo = storage.NewRepository(db)

	// Load saved settings, but CLI flag overrides saved port
	saved := settings.Load()
	if a.port == "" {
		a.port = saved.Port
	}
	if a.port == "" {
		a.port = "2525"
	}
	addr := fmt.Sprintf("127.0.0.1:%s", a.port)

	a.smtpServer, err = smtp.NewServer(addr, a.handleEmail)
	if err != nil {
		log.Fatalf("Failed to start SMTP server on %s: %v\nIs another process using port %s?", addr, err, a.port)
	}
	go a.smtpServer.Start()

	log.Printf("MailTap ready — SMTP on %s, DB at %s", addr, dbPath)
}

func (a *App) GetPort() string {
	return a.port
}

func (a *App) GetSettings() settings.Settings {
	s := settings.Load()
	s.Port = a.port
	return s
}

func (a *App) SaveSettings(s settings.Settings) error {
	return settings.Save(s)
}

func (a *App) shutdown(ctx context.Context) {
	if a.smtpServer != nil {
		a.smtpServer.Stop()
	}
}

func (a *App) handleEmail(email *models.Email) {
	if err := a.repo.InsertEmail(email); err != nil {
		log.Printf("Failed to store email: %v", err)
		return
	}

	wailsruntime.EventsEmit(a.ctx, "mailtap:email-received", email.ToSummary())

	cfg := settings.Load()

	if cfg.Notifications {
		go notify.Send("MailTap", fmt.Sprintf("From: %s\n%s", email.From, email.Subject))
	}

	if cfg.AutoCopyCodes {
		if code, ok := otp.Detect(email); ok {
			go func() {
				if err := wailsruntime.ClipboardSetText(a.ctx, code); err != nil {
					log.Printf("Failed to copy OTP code to clipboard: %v", err)
				} else {
					notify.Send("MailTap", fmt.Sprintf("Code %s copied to clipboard", code))
				}
			}()
		}
	}
}

func (a *App) ListEmails(search string, offset int, limit int) ([]models.EmailSummary, error) {
	return a.repo.ListEmails(search, offset, limit)
}

func (a *App) GetEmail(id string) (*models.Email, error) {
	return a.repo.GetEmail(id)
}

func (a *App) DeleteEmail(id string) error {
	err := a.repo.DeleteEmail(id)
	if err == nil {
		wailsruntime.EventsEmit(a.ctx, "mailtap:email-deleted", id)
	}
	return err
}

func (a *App) DeleteAllEmails() error {
	err := a.repo.DeleteAllEmails()
	if err == nil {
		wailsruntime.EventsEmit(a.ctx, "mailtap:emails-cleared")
	}
	return err
}

func (a *App) MarkAsRead(id string) error {
	return a.repo.MarkAsRead(id)
}

func (a *App) GetEmailCount() (int, error) {
	return a.repo.GetEmailCount()
}

func (a *App) GetRawSource(id string) (string, error) {
	email, err := a.repo.GetEmail(id)
	if err != nil {
		return "", err
	}
	return string(email.RawMessage), nil
}

func (a *App) SaveAttachment(attachmentID string) (string, error) {
	att, err := a.repo.GetAttachment(attachmentID)
	if err != nil {
		return "", err
	}

	savePath, err := wailsruntime.SaveFileDialog(a.ctx, wailsruntime.SaveDialogOptions{
		DefaultFilename: att.Filename,
		Title:           "Save Attachment",
	})
	if err != nil || savePath == "" {
		return "", err
	}

	dir := filepath.Dir(savePath)
	os.MkdirAll(dir, 0755)
	if err := os.WriteFile(savePath, att.Data, 0644); err != nil {
		return "", err
	}
	return savePath, nil
}
