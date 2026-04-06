package main

import (
	"fmt"
	"net/smtp"
	"testing"
	"time"

	"mailtap/internal/models"
	smtpserver "mailtap/internal/smtp"
	"mailtap/internal/storage"
)

func TestFullEmailFlow(t *testing.T) {
	// Set up in-memory DB
	db, err := storage.NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()
	repo := storage.NewRepository(db)

	// Handler that stores emails
	handler := func(email *models.Email) {
		if err := repo.InsertEmail(email); err != nil {
			t.Errorf("failed to insert email: %v", err)
		}
	}

	// Start SMTP server on random port
	srv, err := smtpserver.NewServer("127.0.0.1:0", handler)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	go srv.Start()
	defer srv.Stop()
	time.Sleep(100 * time.Millisecond)

	// Send a test email
	msg := "From: test@example.com\r\nTo: dev@localhost\r\nSubject: Integration Test\r\n" +
		"Content-Type: text/plain\r\n\r\nHello from integration test"
	err = smtp.SendMail(srv.Addr(), nil, "test@example.com", []string{"dev@localhost"}, []byte(msg))
	if err != nil {
		t.Fatalf("failed to send email: %v", err)
	}
	time.Sleep(300 * time.Millisecond)

	// Verify email was stored
	emails, err := repo.ListEmails("", 0, 10)
	if err != nil {
		t.Fatalf("failed to list emails: %v", err)
	}
	if len(emails) != 1 {
		t.Fatalf("expected 1 email, got %d", len(emails))
	}
	if emails[0].Subject != "Integration Test" {
		t.Errorf("wrong subject: %s", emails[0].Subject)
	}

	// Verify full email retrieval
	full, err := repo.GetEmail(emails[0].ID)
	if err != nil {
		t.Fatalf("failed to get email: %v", err)
	}
	if full.From != "test@example.com" {
		t.Errorf("wrong from: %s", full.From)
	}
	if full.PlainBody == "" {
		t.Error("plain body should not be empty")
	}

	// Verify email count
	count, err := repo.GetEmailCount()
	if err != nil {
		t.Fatalf("failed to get count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}

	// Test delete
	err = repo.DeleteEmail(full.ID)
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}
	count, _ = repo.GetEmailCount()
	if count != 0 {
		t.Errorf("expected count 0 after delete, got %d", count)
	}
}

func TestMultipleEmailFlow(t *testing.T) {
	db, _ := storage.NewDB(":memory:")
	defer db.Close()
	repo := storage.NewRepository(db)

	handler := func(email *models.Email) {
		repo.InsertEmail(email)
	}

	srv, _ := smtpserver.NewServer("127.0.0.1:0", handler)
	go srv.Start()
	defer srv.Stop()
	time.Sleep(100 * time.Millisecond)

	// Send 3 emails
	for i := 0; i < 3; i++ {
		msg := fmt.Sprintf("From: sender%d@test.com\r\nTo: dev@localhost\r\nSubject: Email %d\r\nContent-Type: text/plain\r\n\r\nBody %d", i, i, i)
		smtp.SendMail(srv.Addr(), nil, fmt.Sprintf("sender%d@test.com", i), []string{"dev@localhost"}, []byte(msg))
	}
	time.Sleep(500 * time.Millisecond)

	emails, _ := repo.ListEmails("", 0, 10)
	if len(emails) != 3 {
		t.Fatalf("expected 3 emails, got %d", len(emails))
	}

	// Test search
	results, _ := repo.ListEmails("Email 1", 0, 10)
	if len(results) != 1 {
		t.Errorf("expected 1 search result, got %d", len(results))
	}

	// Test delete all
	repo.DeleteAllEmails()
	count, _ := repo.GetEmailCount()
	if count != 0 {
		t.Errorf("expected 0 after delete all, got %d", count)
	}
}
