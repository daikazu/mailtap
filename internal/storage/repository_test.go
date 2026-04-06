package storage

import (
	"fmt"
	"testing"
	"time"

	"mailtap/internal/models"
)

func setupTestDB(t *testing.T) *Repository {
	t.Helper()
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	repo := NewRepository(db)
	t.Cleanup(func() { db.Close() })
	return repo
}

func TestInsertAndGetEmail(t *testing.T) {
	repo := setupTestDB(t)
	email := &models.Email{
		ID: "test-1", From: "sender@test.com", To: []string{"recipient@test.com"},
		CC: []string{}, BCC: []string{}, Subject: "Test Email",
		HTMLBody: "<h1>Hello</h1>", PlainBody: "Hello",
		RawMessage: []byte("raw message content"),
		Headers:    map[string][]string{"From": {"sender@test.com"}},
		ReceivedAt: time.Now(),
	}
	err := repo.InsertEmail(email)
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	got, err := repo.GetEmail("test-1")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Subject != "Test Email" {
		t.Errorf("expected subject 'Test Email', got '%s'", got.Subject)
	}
	if got.From != "sender@test.com" {
		t.Errorf("expected from 'sender@test.com', got '%s'", got.From)
	}
}

func TestListEmails(t *testing.T) {
	repo := setupTestDB(t)
	for i := 0; i < 5; i++ {
		repo.InsertEmail(&models.Email{
			ID: fmt.Sprintf("test-%d", i), From: "sender@test.com", To: []string{"r@test.com"},
			Subject: fmt.Sprintf("Email %d", i), PlainBody: "body",
			ReceivedAt: time.Now().Add(time.Duration(i) * time.Minute),
		})
	}
	results, err := repo.ListEmails("", 0, 10)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(results) != 5 {
		t.Errorf("expected 5, got %d", len(results))
	}
	if results[0].Subject != "Email 4" {
		t.Errorf("expected newest first, got '%s'", results[0].Subject)
	}
}

func TestSearchEmails(t *testing.T) {
	repo := setupTestDB(t)
	repo.InsertEmail(&models.Email{ID: "a", From: "alice@test.com", Subject: "Hello World", PlainBody: "body", ReceivedAt: time.Now(), To: []string{"r@test.com"}})
	repo.InsertEmail(&models.Email{ID: "b", From: "bob@test.com", Subject: "Goodbye", PlainBody: "body", ReceivedAt: time.Now(), To: []string{"r@test.com"}})
	results, err := repo.ListEmails("Hello", 0, 10)
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1, got %d", len(results))
	}
}

func TestDeleteEmail(t *testing.T) {
	repo := setupTestDB(t)
	repo.InsertEmail(&models.Email{ID: "del-1", From: "s@t.com", Subject: "Delete Me", ReceivedAt: time.Now(), To: []string{"r@t.com"}})
	err := repo.DeleteEmail("del-1")
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	_, err = repo.GetEmail("del-1")
	if err == nil {
		t.Error("expected error getting deleted email")
	}
}

func TestDeleteAllEmails(t *testing.T) {
	repo := setupTestDB(t)
	for i := 0; i < 3; i++ {
		repo.InsertEmail(&models.Email{ID: fmt.Sprintf("d-%d", i), From: "s@t.com", Subject: "X", ReceivedAt: time.Now(), To: []string{"r@t.com"}})
	}
	err := repo.DeleteAllEmails()
	if err != nil {
		t.Fatalf("delete all failed: %v", err)
	}
	count, _ := repo.GetEmailCount()
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}

func TestMarkAsRead(t *testing.T) {
	repo := setupTestDB(t)
	repo.InsertEmail(&models.Email{ID: "read-1", From: "s@t.com", Subject: "Read Me", ReceivedAt: time.Now(), To: []string{"r@t.com"}})
	err := repo.MarkAsRead("read-1")
	if err != nil {
		t.Fatalf("mark read failed: %v", err)
	}
	email, _ := repo.GetEmail("read-1")
	if !email.Read {
		t.Error("expected email to be marked as read")
	}
}

func TestAttachments(t *testing.T) {
	repo := setupTestDB(t)
	email := &models.Email{
		ID: "att-1", From: "s@t.com", Subject: "With Attachment", ReceivedAt: time.Now(), To: []string{"r@t.com"},
		Attachments: []models.Attachment{{ID: "file-1", Filename: "test.pdf", ContentType: "application/pdf", Size: 1024, Data: []byte("fake pdf data")}},
	}
	repo.InsertEmail(email)
	got, _ := repo.GetEmail("att-1")
	if len(got.Attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(got.Attachments))
	}
	if got.Attachments[0].Filename != "test.pdf" {
		t.Errorf("expected 'test.pdf', got '%s'", got.Attachments[0].Filename)
	}
}
