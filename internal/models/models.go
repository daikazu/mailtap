package models

import "time"

type Email struct {
	ID          string              `json:"id"`
	From        string              `json:"from"`
	To          []string            `json:"to"`
	CC          []string            `json:"cc"`
	BCC         []string            `json:"bcc"`
	Subject     string              `json:"subject"`
	HTMLBody    string              `json:"htmlBody"`
	PlainBody   string              `json:"plainBody"`
	RawMessage  []byte              `json:"-"`
	Headers     map[string][]string `json:"headers"`
	Attachments []Attachment        `json:"attachments"`
	ReceivedAt  time.Time           `json:"receivedAt"`
	Read        bool                `json:"read"`
}

type EmailSummary struct {
	ID              string    `json:"id"`
	From            string    `json:"from"`
	To              []string  `json:"to"`
	Subject         string    `json:"subject"`
	Preview         string    `json:"preview"`
	AttachmentCount int       `json:"attachmentCount"`
	ReceivedAt      time.Time `json:"receivedAt"`
	Read            bool      `json:"read"`
}

func (e *Email) ToSummary() EmailSummary {
	preview := e.PlainBody
	if len(preview) > 100 {
		preview = preview[:100]
	}
	return EmailSummary{
		ID:              e.ID,
		From:            e.From,
		To:              e.To,
		Subject:         e.Subject,
		Preview:         preview,
		AttachmentCount: len(e.Attachments),
		ReceivedAt:      e.ReceivedAt,
		Read:            e.Read,
	}
}

type Attachment struct {
	ID          string `json:"id"`
	EmailID     string `json:"emailId"`
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
	Data        []byte `json:"-"`
}

