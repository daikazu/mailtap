package smtp

import (
	"net/smtp"
	"strings"
	"testing"
	"time"

	"mailtap/internal/models"
)

func TestSMTPServerReceivesEmail(t *testing.T) {
	var received *models.Email
	handler := func(email *models.Email) { received = email }

	srv, err := NewServer("127.0.0.1:0", handler)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	go srv.Start()
	defer srv.Stop()
	time.Sleep(100 * time.Millisecond)

	addr := srv.Addr()
	msg := "From: sender@test.com\r\nTo: recipient@test.com\r\nSubject: Test\r\n\r\nHello"
	err = smtp.SendMail(addr, nil, "sender@test.com", []string{"recipient@test.com"}, []byte(msg))
	if err != nil {
		t.Fatalf("send failed: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	if received == nil {
		t.Fatal("expected to receive an email")
	}
	if !strings.Contains(received.Subject, "Test") {
		t.Errorf("expected subject 'Test', got '%s'", received.Subject)
	}
}
