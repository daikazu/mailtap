# OTP Code Auto-Copy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Detect a one-time code in a received email and automatically copy it to the clipboard plus fire a native notification, controlled by a default-on setting.

**Architecture:** A new pure-Go `internal/otp` package exposes `Detect(email) (code, found)` with no Wails dependency. `app.go:handleEmail()` calls it after storing the email and, when the new `AutoCopyCodes` setting is enabled, writes the code via Wails' cross-platform clipboard API and sends a notification through the existing `internal/notify` package. A toggle is added to the Settings UI.

**Tech Stack:** Go (stdlib `regexp`/`strings`, `modernc.org/sqlite` unaffected), Wails v2 runtime (`ClipboardSetText`), Svelte 3 frontend, Go stdlib `testing` (no testify — match existing test style).

**Spec:** `docs/superpowers/specs/2026-05-15-otp-code-auto-copy-design.md`

---

## File Structure

- **Create** `internal/otp/otp.go` — detection logic. Single responsibility: turn an `*models.Email` into a detected code or nothing. No Wails/clipboard/notify imports.
- **Create** `internal/otp/otp_test.go` — table-driven detection tests (positives + negatives).
- **Modify** `internal/settings/settings.go` — add `AutoCopyCodes` field and default.
- **Create** `internal/settings/settings_test.go` — assert default + legacy-file backward compatibility.
- **Modify** `app.go` — wire `otp.Detect` + clipboard + notify into `handleEmail()`.
- **Modify** `frontend/src/components/Settings.svelte` — add the toggle and include it in the save payload.

---

## Task 1: `internal/otp` detection package

**Files:**
- Create: `internal/otp/otp.go`
- Test: `internal/otp/otp_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/otp/otp_test.go`:

```go
package otp

import (
	"testing"

	"mailtap/internal/models"
)

func TestDetect(t *testing.T) {
	cases := []struct {
		name      string
		subject   string
		plainBody string
		htmlBody  string
		wantCode  string
		wantFound bool
	}{
		{
			name:      "github numeric code",
			subject:   "Your GitHub verification code",
			plainBody: "Verification code: 123456\nThis code expires in 10 minutes.",
			wantCode:  "123456",
			wantFound: true,
		},
		{
			name:      "google inline code",
			subject:   "Security alert",
			plainBody: "Your Google verification code is 482913. Do not share it.",
			wantCode:  "482913",
			wantFound: true,
		},
		{
			name:      "grouped code normalized",
			subject:   "One-time passcode",
			plainBody: "Your code is 123-456 and is valid for 5 minutes.",
			wantCode:  "123456",
			wantFound: true,
		},
		{
			name:      "alphanumeric code",
			subject:   "Confirm your email",
			plainBody: "Use code A1B2C3 to verify your account.",
			wantCode:  "A1B2C3",
			wantFound: true,
		},
		{
			name:      "html fallback when plain empty",
			subject:   "Verify",
			htmlBody:  "<p>Your verification code is <b>778812</b></p>",
			wantCode:  "778812",
			wantFound: true,
		},
		{
			name:      "no keyword - order number ignored",
			subject:   "Your receipt",
			plainBody: "Total: $123.45. Order number 88231 ships Monday.",
			wantFound: false,
		},
		{
			name:      "currency excluded even near keyword",
			subject:   "Payment confirmation",
			plainBody: "We confirm your payment of $1234 was received.",
			wantFound: false,
		},
		{
			name:      "year excluded",
			subject:   "Security update",
			plainBody: "Your security settings were reviewed in 2025.",
			wantFound: false,
		},
		{
			name:      "keyword but no token",
			subject:   "Action needed",
			plainBody: "Enter your verification code on the website to continue.",
			wantFound: false,
		},
		{
			name:      "tracking number without keyword",
			subject:   "Shipped",
			plainBody: "Your package AB1234567 is on its way via the carrier.",
			wantFound: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			email := &models.Email{
				Subject:   tc.subject,
				PlainBody: tc.plainBody,
				HTMLBody:  tc.htmlBody,
			}
			code, found := Detect(email)
			if found != tc.wantFound {
				t.Fatalf("found = %v, want %v (code=%q)", found, tc.wantFound, code)
			}
			if found && code != tc.wantCode {
				t.Errorf("code = %q, want %q", code, tc.wantCode)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/otp/ -v`
Expected: build failure — `otp.go` does not exist / `undefined: Detect`.

- [ ] **Step 3: Write minimal implementation**

Create `internal/otp/otp.go`:

```go
// Package otp detects one-time codes (verification / 2FA codes) in emails.
package otp

import (
	"regexp"
	"strings"

	"mailtap/internal/models"
)

// keywords are lower-cased context words that must appear near a candidate
// token for it to be treated as a one-time code.
var keywords = []string{
	"verification", "verify", "one-time", "one time", "passcode",
	"otp", "2fa", "two-factor", "two factor", "security code",
	"confirmation", "confirm", "authentication", "authenticate",
	"access code", "login code", "sign-in", "sign in",
	"security", "pin", "code",
}

var (
	reDigits  = regexp.MustCompile(`\b\d{4,8}\b`)
	reGrouped = regexp.MustCompile(`\b\d{3}[- ]\d{3}\b`)
	reAlnum   = regexp.MustCompile(`\b[A-Z0-9]{6,8}\b`)
	reTag     = regexp.MustCompile(`<[^>]+>`)
	reYear    = regexp.MustCompile(`^(?:19|20)\d{2}$`)
)

const proximityWindow = 50

type candidate struct {
	raw   string // exact substring as it appears in text
	value string // normalized code (grouping separators stripped)
	pos   int    // byte offset of raw in text
}

// Detect returns the one-time code found in the email and whether one was
// found. It only returns a code when a candidate token appears within
// proximityWindow characters of a context keyword. Ambiguous or absent
// matches return ("", false) so the clipboard is never touched on uncertainty.
func Detect(email *models.Email) (string, bool) {
	text := email.Subject + "\n" + email.PlainBody
	if strings.TrimSpace(email.PlainBody) == "" && email.HTMLBody != "" {
		text = email.Subject + "\n" + stripHTML(email.HTMLBody)
	}

	lower := strings.ToLower(text)
	kwPos := keywordPositions(lower)
	if len(kwPos) == 0 {
		return "", false
	}

	cands := collectCandidates(text)

	best := ""
	bestDist := proximityWindow + 1
	for _, c := range cands {
		if isExcluded(c, text) {
			continue
		}
		d := nearestDistance(c.pos, len(c.raw), kwPos)
		if d <= proximityWindow && d < bestDist {
			bestDist = d
			best = c.value
		}
	}
	if best == "" {
		return "", false
	}
	return best, true
}

func keywordPositions(lower string) []int {
	var positions []int
	for _, kw := range keywords {
		from := 0
		for {
			i := strings.Index(lower[from:], kw)
			if i < 0 {
				break
			}
			positions = append(positions, from+i)
			from += i + len(kw)
		}
	}
	return positions
}

func collectCandidates(text string) []candidate {
	var cands []candidate

	for _, m := range reGrouped.FindAllStringIndex(text, -1) {
		raw := text[m[0]:m[1]]
		norm := strings.NewReplacer("-", "", " ", "").Replace(raw)
		cands = append(cands, candidate{raw: raw, value: norm, pos: m[0]})
	}
	for _, m := range reDigits.FindAllStringIndex(text, -1) {
		raw := text[m[0]:m[1]]
		cands = append(cands, candidate{raw: raw, value: raw, pos: m[0]})
	}
	for _, m := range reAlnum.FindAllStringIndex(text, -1) {
		raw := text[m[0]:m[1]]
		if !strings.ContainsAny(raw, "0123456789") {
			continue // plain uppercase word, not a code
		}
		cands = append(cands, candidate{raw: raw, value: raw, pos: m[0]})
	}
	return cands
}

func isExcluded(c candidate, text string) bool {
	if len(c.value) == 4 && reYear.MatchString(c.value) {
		return true
	}
	if c.pos > 0 && text[c.pos-1] == '$' {
		return true
	}
	end := c.pos + len(c.raw)
	if end+2 < len(text) && text[end] == '.' && isDigit(text[end+1]) && isDigit(text[end+2]) {
		return true
	}
	return false
}

func nearestDistance(pos, length int, kw []int) int {
	best := 1 << 30
	for _, k := range kw {
		var d int
		switch {
		case k < pos:
			d = pos - k
		case k > pos+length:
			d = k - (pos + length)
		default:
			d = 0
		}
		if d < best {
			best = d
		}
	}
	return best
}

func stripHTML(s string) string {
	return reTag.ReplaceAllString(s, " ")
}

func isDigit(b byte) bool { return b >= '0' && b <= '9' }
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/otp/ -v`
Expected: PASS — all subtests of `TestDetect` pass.

- [ ] **Step 5: Commit**

```bash
git add internal/otp/otp.go internal/otp/otp_test.go
git commit -m "feat: add otp package for one-time code detection

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 2: `AutoCopyCodes` setting

**Files:**
- Modify: `internal/settings/settings.go`
- Test: `internal/settings/settings_test.go`

- [ ] **Step 1: Write the failing test**

Create `internal/settings/settings_test.go`:

```go
package settings

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultAutoCopyCodesIsTrue(t *testing.T) {
	if !Default().AutoCopyCodes {
		t.Errorf("Default().AutoCopyCodes = false, want true")
	}
}

func TestLoadLegacyFilePreservesAutoCopyDefault(t *testing.T) {
	dir := t.TempDir()
	legacy := filepath.Join(dir, "settings.json")
	// A settings file written before this feature existed (no autoCopyCodes key).
	if err := os.WriteFile(legacy, []byte(`{"port":"2525","notifications":true,"theme":"dark"}`), 0644); err != nil {
		t.Fatalf("write legacy file: %v", err)
	}

	mu.Lock()
	orig := filePath
	filePath = legacy
	mu.Unlock()
	defer func() {
		mu.Lock()
		filePath = orig
		mu.Unlock()
	}()

	s := Load()
	if !s.AutoCopyCodes {
		t.Errorf("legacy load AutoCopyCodes = false, want true (missing key must keep default)")
	}
	if s.Theme != "dark" {
		t.Errorf("legacy load Theme = %q, want \"dark\" (sanity: file was read)", s.Theme)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/settings/ -v`
Expected: build failure — `s.AutoCopyCodes` undefined (field does not exist yet).

- [ ] **Step 3: Add the field and default**

In `internal/settings/settings.go`, change the `Settings` struct from:

```go
type Settings struct {
	Port          string `json:"port"`
	Notifications bool   `json:"notifications"`
	Theme         string `json:"theme"`
}
```

to:

```go
type Settings struct {
	Port          string `json:"port"`
	Notifications bool   `json:"notifications"`
	Theme         string `json:"theme"`
	AutoCopyCodes bool   `json:"autoCopyCodes"`
}
```

And change `Default()` from:

```go
func Default() Settings {
	return Settings{
		Port:          "2525",
		Notifications: true,
		Theme:         "system",
	}
}
```

to:

```go
func Default() Settings {
	return Settings{
		Port:          "2525",
		Notifications: true,
		Theme:         "system",
		AutoCopyCodes: true,
	}
}
```

(No change to `Load()` is needed: it starts from `Default()` then unmarshals the
file over it, so a legacy file lacking `autoCopyCodes` keeps the `true` default.)

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/settings/ -v`
Expected: PASS — both tests pass.

- [ ] **Step 5: Commit**

```bash
git add internal/settings/settings.go internal/settings/settings_test.go
git commit -m "feat: add AutoCopyCodes setting (default on)

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 3: Wire detection into `handleEmail`

**Files:**
- Modify: `app.go:79-90` (`handleEmail`)

This is glue between already-tested units (`otp.Detect`, `settings.Load`, the
Wails clipboard runtime, and `notify.Send`). It is verified by build + vet here
and by the manual end-to-end check below; the detection logic itself is covered
by Task 1's tests.

- [ ] **Step 1: Add the `otp` import**

In `app.go`, the import block currently ends:

```go
	"mailtap/internal/models"
	"mailtap/internal/notify"
	"mailtap/internal/settings"
	"mailtap/internal/smtp"
	"mailtap/internal/storage"
)
```

Add the `otp` import (keep alphabetical order):

```go
	"mailtap/internal/models"
	"mailtap/internal/notify"
	"mailtap/internal/otp"
	"mailtap/internal/settings"
	"mailtap/internal/smtp"
	"mailtap/internal/storage"
)
```

- [ ] **Step 2: Add the auto-copy block in `handleEmail`**

`handleEmail` currently ends:

```go
	wailsruntime.EventsEmit(a.ctx, "mailtap:email-received", email.ToSummary())

	if settings.Load().Notifications {
		go notify.Send("MailTap", fmt.Sprintf("From: %s\n%s", email.From, email.Subject))
	}
}
```

Replace that with:

```go
	wailsruntime.EventsEmit(a.ctx, "mailtap:email-received", email.ToSummary())

	cfg := settings.Load()

	if cfg.Notifications {
		go notify.Send("MailTap", fmt.Sprintf("From: %s\n%s", email.From, email.Subject))
	}

	if cfg.AutoCopyCodes {
		if code, ok := otp.Detect(email); ok {
			if err := wailsruntime.ClipboardSetText(a.ctx, code); err != nil {
				log.Printf("Failed to copy OTP code to clipboard: %v", err)
			} else {
				go notify.Send("MailTap", fmt.Sprintf("Code %s copied to clipboard", code))
			}
		}
	}
}
```

(`settings.Load()` is now called once and reused; behavior of the existing
notification path is unchanged.)

- [ ] **Step 3: Build and vet**

Run: `go build ./... && go vet ./...`
Expected: no output, exit code 0 (compiles cleanly, no vet warnings).

- [ ] **Step 4: Run the full Go test suite**

Run: `go test ./...`
Expected: PASS for all packages (`otp`, `settings`, `parser`, `storage`, `smtp`, and the root integration test).

- [ ] **Step 5: Manual end-to-end check**

Start the app: `wails dev` (in a second terminal). Send a test email with a code:

```bash
printf 'From: auth@example.com\r\nTo: me@example.com\r\nSubject: Your verification code\r\n\r\nYour verification code is 482913. It expires in 10 minutes.\r\n' \
  | nc 127.0.0.1 2525
```

Expected: a native notification "Code 482913 copied to clipboard" appears, and
pasting (Cmd/Ctrl+V) anywhere yields `482913`. Then send an email with no code
keyword and confirm the clipboard is **not** changed.

- [ ] **Step 6: Commit**

```bash
git add app.go
git commit -m "feat: auto-copy detected OTP code to clipboard on receive

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Task 4: Settings UI toggle

**Files:**
- Modify: `frontend/src/components/Settings.svelte`

- [ ] **Step 1: Add the state variable**

In the `<script>` block, the declarations are:

```js
  let port = '';
  let notifications = true;
  let theme = 'system';
  let needsRestart = false;
  let saved = false;
```

Add `autoCopyCodes`:

```js
  let port = '';
  let notifications = true;
  let autoCopyCodes = true;
  let theme = 'system';
  let needsRestart = false;
  let saved = false;
```

- [ ] **Step 2: Load it in `onMount`**

`onMount` currently is:

```js
  onMount(async () => {
    const s = await api.getSettings();
    port = s.port;
    notifications = s.notifications;
    theme = s.theme;
  });
```

Change it to:

```js
  onMount(async () => {
    const s = await api.getSettings();
    port = s.port;
    notifications = s.notifications;
    autoCopyCodes = s.autoCopyCodes;
    theme = s.theme;
  });
```

- [ ] **Step 3: Include it in the save payload**

In `save()`, this line:

```js
    await api.saveSettings({ port, notifications, theme });
```

becomes:

```js
    await api.saveSettings({ port, notifications, autoCopyCodes, theme });
```

- [ ] **Step 4: Add the toggle row in the template**

After the Notifications `<label class="setting-row">…</label>` block (the one
ending right before the Theme row), insert a new row:

```svelte
      <label class="setting-row">
        <span class="setting-label">Auto-copy detected codes</span>
        <button
          class="toggle"
          class:active={autoCopyCodes}
          on:click={() => autoCopyCodes = !autoCopyCodes}
        >
          <span class="toggle-knob"></span>
        </button>
      </label>
```

So the order in `.settings-body` is: SMTP Port, Notifications, Auto-copy
detected codes, Theme. (No CSS changes — it reuses the existing `.setting-row`
and `.toggle` styles.)

- [ ] **Step 5: Build the frontend / regenerate bindings**

Run: `wails build` (this regenerates the Wails TS bindings for the new
`Settings.AutoCopyCodes` field and compiles the frontend).
Expected: build succeeds with no errors.

(Note: `api.saveSettings` forwards a plain object and `api.getSettings` returns
parsed JSON, so no change to `frontend/src/lib/api.js` is required; the
regenerated bindings just keep the TypeScript types in sync.)

- [ ] **Step 6: Manual UI check**

Run `wails dev`, open Settings (gear / command palette). Confirm the
"Auto-copy detected codes" toggle appears between Notifications and Theme,
defaults to on, and that toggling it off + Save + reopening Settings shows it
still off. With it off, sending the Task 3 test email must NOT change the
clipboard; with it on, it must.

- [ ] **Step 7: Commit**

```bash
git add frontend/src/components/Settings.svelte
git commit -m "feat: add auto-copy codes toggle to settings UI

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>"
```

---

## Final Verification

- [ ] Run `go test ./...` — all packages PASS.
- [ ] Run `go vet ./...` — no warnings.
- [ ] Run `wails build` — succeeds.
- [ ] Manual: code email → notification + clipboard set; no-keyword email → clipboard untouched; toggle off → no auto-copy.
