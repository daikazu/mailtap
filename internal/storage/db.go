package storage

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func NewDB(dbPath string) (*sql.DB, error) {
	if dbPath != ":memory:" {
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA foreign_keys=ON")

	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS emails (
		id TEXT PRIMARY KEY,
		from_addr TEXT NOT NULL,
		to_addrs TEXT NOT NULL,
		cc_addrs TEXT NOT NULL DEFAULT '',
		bcc_addrs TEXT NOT NULL DEFAULT '',
		subject TEXT NOT NULL DEFAULT '',
		html_body TEXT NOT NULL DEFAULT '',
		plain_body TEXT NOT NULL DEFAULT '',
		raw_message BLOB,
		headers TEXT NOT NULL DEFAULT '{}',
		received_at DATETIME NOT NULL,
		read BOOLEAN NOT NULL DEFAULT 0
	);
	CREATE TABLE IF NOT EXISTS attachments (
		id TEXT PRIMARY KEY,
		email_id TEXT NOT NULL,
		filename TEXT NOT NULL,
		content_type TEXT NOT NULL,
		size INTEGER NOT NULL,
		data BLOB NOT NULL,
		FOREIGN KEY (email_id) REFERENCES emails(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_emails_received_at ON emails(received_at);
	CREATE INDEX IF NOT EXISTS idx_emails_from ON emails(from_addr);
	CREATE INDEX IF NOT EXISTS idx_emails_subject ON emails(subject);
	CREATE INDEX IF NOT EXISTS idx_attachments_email_id ON attachments(email_id);
	`
	_, err := db.Exec(schema)
	return err
}

func DefaultDBPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mailtap", "mailtap.db")
}
