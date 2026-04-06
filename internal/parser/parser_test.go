package parser

import (
	"strings"
	"testing"
)

func TestParseSimpleEmail(t *testing.T) {
	raw := "From: sender@test.com\r\nTo: recipient@test.com\r\nSubject: Test\r\nContent-Type: text/plain\r\n\r\nHello World"
	email, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if email.From != "sender@test.com" {
		t.Errorf("expected from 'sender@test.com', got '%s'", email.From)
	}
	if email.Subject != "Test" {
		t.Errorf("expected subject 'Test', got '%s'", email.Subject)
	}
	if !strings.Contains(email.PlainBody, "Hello World") {
		t.Errorf("expected plain body to contain 'Hello World', got '%s'", email.PlainBody)
	}
}

func TestParseHTMLEmail(t *testing.T) {
	raw := "From: sender@test.com\r\nTo: recipient@test.com\r\nSubject: HTML Test\r\n" +
		"Content-Type: multipart/alternative; boundary=boundary\r\n\r\n" +
		"--boundary\r\nContent-Type: text/plain\r\n\r\nPlain text\r\n" +
		"--boundary\r\nContent-Type: text/html\r\n\r\n<h1>HTML</h1>\r\n" +
		"--boundary--"
	email, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if email.HTMLBody != "<h1>HTML</h1>" {
		t.Errorf("expected HTML body '<h1>HTML</h1>', got '%s'", email.HTMLBody)
	}
	if email.PlainBody != "Plain text" {
		t.Errorf("expected plain body 'Plain text', got '%s'", email.PlainBody)
	}
}

func TestParseMalformedEmail(t *testing.T) {
	raw := "this is not a valid email"
	email, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("should not error on malformed: %v", err)
	}
	if email.ID == "" {
		t.Error("should still have an ID assigned")
	}
}
