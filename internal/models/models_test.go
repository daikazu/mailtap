package models

import (
	"testing"
	"time"
)

func TestEmailSummary(t *testing.T) {
	e := Email{
		ID:        "test-uuid",
		From:      "sender@example.com",
		To:        []string{"recipient@example.com"},
		Subject:   "Test Subject",
		PlainBody:  "Hello world, this is a test email body with enough text",
		ReceivedAt: time.Now(),
	}
	summary := e.ToSummary()
	if summary.ID != e.ID {
		t.Errorf("expected ID %s, got %s", e.ID, summary.ID)
	}
	if summary.From != e.From {
		t.Errorf("expected From %s, got %s", e.From, summary.From)
	}
	if len(summary.Preview) > 100 {
		t.Error("preview should be truncated to 100 chars")
	}
}

