package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"mailtap/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) InsertEmail(e *models.Email) error {
	headersJSON, _ := json.Marshal(e.Headers)

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO emails (id, from_addr, to_addrs, cc_addrs, bcc_addrs, subject,
		html_body, plain_body, raw_message, headers, received_at, read)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.From, strings.Join(e.To, ","), strings.Join(e.CC, ","), strings.Join(e.BCC, ","),
		e.Subject, e.HTMLBody, e.PlainBody, e.RawMessage, string(headersJSON),
		e.ReceivedAt, e.Read)
	if err != nil {
		return err
	}

	for _, att := range e.Attachments {
		_, err = tx.Exec(`INSERT INTO attachments (id, email_id, filename, content_type, size, data) VALUES (?, ?, ?, ?, ?, ?)`,
			att.ID, e.ID, att.Filename, att.ContentType, att.Size, att.Data)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) GetEmail(id string) (*models.Email, error) {
	e := &models.Email{}
	var toStr, ccStr, bccStr, headersStr string

	err := r.db.QueryRow(`SELECT id, from_addr, to_addrs, cc_addrs, bcc_addrs, subject,
		html_body, plain_body, raw_message, headers, received_at, read
		FROM emails WHERE id = ?`, id).Scan(
		&e.ID, &e.From, &toStr, &ccStr, &bccStr, &e.Subject, &e.HTMLBody, &e.PlainBody,
		&e.RawMessage, &headersStr, &e.ReceivedAt, &e.Read)
	if err != nil {
		return nil, err
	}

	e.To = splitNonEmpty(toStr, ",")
	e.CC = splitNonEmpty(ccStr, ",")
	e.BCC = splitNonEmpty(bccStr, ",")
	json.Unmarshal([]byte(headersStr), &e.Headers)

	rows, err := r.db.Query(`SELECT id, filename, content_type, size, data FROM attachments WHERE email_id = ?`, id)
	if err != nil {
		return e, nil
	}
	defer rows.Close()

	for rows.Next() {
		var att models.Attachment
		rows.Scan(&att.ID, &att.Filename, &att.ContentType, &att.Size, &att.Data)
		att.EmailID = id
		e.Attachments = append(e.Attachments, att)
	}

	return e, nil
}

func (r *Repository) ListEmails(search string, offset, limit int) ([]models.EmailSummary, error) {
	var rows *sql.Rows
	var err error

	if search == "" {
		rows, err = r.db.Query(`SELECT id, from_addr, to_addrs, subject, plain_body, received_at, read,
			(SELECT COUNT(*) FROM attachments WHERE email_id = emails.id) as att_count
			FROM emails ORDER BY received_at DESC LIMIT ? OFFSET ?`, limit, offset)
	} else {
		like := "%" + search + "%"
		rows, err = r.db.Query(`SELECT id, from_addr, to_addrs, subject, plain_body, received_at, read,
			(SELECT COUNT(*) FROM attachments WHERE email_id = emails.id) as att_count
			FROM emails WHERE subject LIKE ? OR from_addr LIKE ? OR plain_body LIKE ?
			ORDER BY received_at DESC LIMIT ? OFFSET ?`, like, like, like, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.EmailSummary
	for rows.Next() {
		var s models.EmailSummary
		var toStr, plainBody string
		rows.Scan(&s.ID, &s.From, &toStr, &s.Subject, &plainBody, &s.ReceivedAt, &s.Read, &s.AttachmentCount)
		s.To = splitNonEmpty(toStr, ",")
		s.Preview = plainBody
		if len(s.Preview) > 100 {
			s.Preview = s.Preview[:100]
		}
		results = append(results, s)
	}

	return results, nil
}

func (r *Repository) DeleteEmail(id string) error {
	result, err := r.db.Exec("DELETE FROM emails WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("email not found: %s", id)
	}
	return nil
}

func (r *Repository) DeleteAllEmails() error {
	_, err := r.db.Exec("DELETE FROM emails")
	return err
}

func (r *Repository) MarkAsRead(id string) error {
	_, err := r.db.Exec("UPDATE emails SET read = 1 WHERE id = ?", id)
	return err
}

func (r *Repository) GetEmailCount() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM emails").Scan(&count)
	return count, err
}

func (r *Repository) GetAttachment(id string) (*models.Attachment, error) {
	att := &models.Attachment{}
	err := r.db.QueryRow(`SELECT id, email_id, filename, content_type, size, data FROM attachments WHERE id = ?`, id).Scan(
		&att.ID, &att.EmailID, &att.Filename, &att.ContentType, &att.Size, &att.Data)
	return att, err
}

func splitNonEmpty(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}
