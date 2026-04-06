package parser

import (
	"bytes"
	"time"

	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"

	"mailtap/internal/models"
)

func Parse(raw []byte) (*models.Email, error) {
	email := &models.Email{
		ID:         uuid.New().String(),
		RawMessage: raw,
		ReceivedAt: time.Now(),
		Headers:    make(map[string][]string),
	}

	env, err := enmime.ReadEnvelope(bytes.NewReader(raw))
	if err != nil {
		return email, nil // Best-effort: store raw even if parsing fails
	}

	email.From = env.GetHeader("From")
	email.Subject = env.GetHeader("Subject")
	email.HTMLBody = env.HTML
	email.PlainBody = env.Text

	if toList, err := env.AddressList("To"); err == nil {
		for _, addr := range toList {
			email.To = append(email.To, addr.Address)
		}
	}
	if ccList, err := env.AddressList("Cc"); err == nil {
		for _, addr := range ccList {
			email.CC = append(email.CC, addr.Address)
		}
	}
	if bccList, err := env.AddressList("Bcc"); err == nil {
		for _, addr := range bccList {
			email.BCC = append(email.BCC, addr.Address)
		}
	}

	for _, key := range env.GetHeaderKeys() {
		vals := env.GetHeaderValues(key)
		email.Headers[key] = vals
	}

	for _, att := range env.Attachments {
		email.Attachments = append(email.Attachments, models.Attachment{
			ID:          uuid.New().String(),
			Filename:    att.FileName,
			ContentType: att.ContentType,
			Size:        int64(len(att.Content)),
			Data:        att.Content,
		})
	}

	return email, nil
}
