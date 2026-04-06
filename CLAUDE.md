# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

MailTap is a local email capture and inspection desktop app built with **Wails v2** (Go backend + Svelte frontend). It runs an SMTP server on `127.0.0.1:2525` (configurable via `-port` flag), captures emails into a SQLite database (`~/.mailtap/mailtap.db`), and displays them in a native desktop window.

## Commands

```bash
# Development (hot-reload for frontend, live Go rebuild)
wails dev

# Production build
wails build

# Run all Go tests
go test ./...

# Run a single package's tests
go test ./internal/storage/

# Run integration tests only
go test -run TestFullEmailFlow

# Frontend dev server (standalone, without Wails)
cd frontend && npm run dev
```

## Architecture

### Go Backend

- **`main.go`** — Entry point. Parses `-port` flag, creates App, runs Wails.
- **`app.go`** — `App` struct is the Wails binding layer. All methods on `App` are callable from the frontend via Wails' auto-generated TypeScript bindings. Handles email lifecycle: receive → store → emit event.
- **`internal/smtp/`** — SMTP server using `go-smtp`. Accepts mail, delegates parsing, calls `EmailHandler` callback. Binds a `net.Listener` at startup (supports `:0` for random port in tests).
- **`internal/parser/`** — Email MIME parsing using `enmime`. Converts raw bytes to `models.Email`.
- **`internal/storage/`** — SQLite persistence via `modernc.org/sqlite` (pure Go, no CGo). `db.go` handles connection/migration, `repository.go` handles CRUD. Uses `:memory:` in tests.
- **`internal/models/`** — Shared data types: `Email`, `EmailSummary`, `Attachment`.

### Frontend (`frontend/`)

Svelte 3 SPA with Vite. No external UI framework — custom CSS with light/dark theme support.

- **`src/lib/api.js`** — Thin wrapper around Wails-generated Go bindings (`wailsjs/go/main/App`).
- **`src/lib/stores.js`** — Svelte stores for app state (email list, selection, theme, search). Contains business logic for email selection, deletion, and list management.
- **`src/lib/hotkeys.js`** — Keyboard shortcut handling.
- **`src/components/`** — UI components: `EmailList`, `DetailPanel`, `TabBar`, `HtmlPreview`, `Headers`, `RawSource`, `Attachments`, `CommandPalette`, etc.

### Data Flow

1. SMTP server receives email → `parser.Parse()` → `models.Email`
2. `App.handleEmail()` stores in SQLite
3. Wails event `mailtap:email-received` emitted to frontend with `EmailSummary`
4. Frontend updates stores reactively; detail views fetch full email on demand

### Key Patterns

- Go ↔ Frontend communication is via Wails bindings (methods on `App` struct) and Wails events (`EventsEmit`/`EventsOn`).
- Frontend events: `mailtap:email-received`, `mailtap:email-deleted`, `mailtap:emails-cleared`.
- Database uses auto-migration on startup — schema is in `storage/db.go`.
- Tests use `:memory:` SQLite and `127.0.0.1:0` for ephemeral SMTP ports.
