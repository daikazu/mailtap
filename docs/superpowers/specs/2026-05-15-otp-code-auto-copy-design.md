# OTP Code Auto-Copy — Design

**Date:** 2026-05-15
**Branch:** `feat-copy-otp`
**Status:** Approved

## Summary

When MailTap receives an email that contains a one-time code (OTP, verification
code, 2FA code, etc.), detect the code and automatically copy it to the system
clipboard, then fire a native notification confirming the copy. This gives the
common developer workflow — "I just triggered a signup/login, paste the code" —
a near-zero-friction path, approximating the macOS Messages experience.

## Background / Constraints

The native macOS/iOS feature (Security Code AutoFill) cannot be used by MailTap:

- It works on the **consuming** side — a text field tagged with the
  `one-time-code` attribute gets a QuickType suggestion. The OS itself scans the
  code **source**, and Apple special-cases its own Messages and Apple Mail apps.
- There is **no public API** for a third-party app to publish/inject a
  freshly-received code into the system suggestion bar.
- macOS Tahoe / iOS 26 expanded scanning to some third-party messaging/email
  clients, but it is still **the OS scanning recognized clients**, not an API an
  app calls. MailTap is a local SMTP catcher (a Wails app), not a registered
  system mail client, and cannot register as a code source.
- The adjacent public API ("Providing one-time passcodes to AutoFill") is for
  credential-provider extensions (password-manager / TOTP apps that own the
  account secret). It requires a separate native app-extension target and only
  surfaces for accounts the app manages — an architectural mismatch for a Wails
  dev tool and not practical here.

**Conclusion:** MailTap must detect and surface the code itself. The chosen
equivalent is auto-copy to the clipboard plus a native notification.

References:

- Apple Developer: Providing one-time passcodes to AutoFill —
  https://developer.apple.com/documentation/authenticationservices/providing-one-time-passcodes-to-autofill
- MacRumors: iOS 26 and macOS Tahoe Expand AutoFill for One-Time Codes —
  https://www.macrumors.com/2025/06/12/ios-26-macos-tahoe-expand-one-time-code-autofill/
- 9to5Mac: Apple's one-time code AutoFill in iOS 26 —
  https://9to5mac.com/2025/08/20/apples-killer-one-time-code-autofill-feature-gets-even-better-in-ios-26/

## Decisions

| Decision | Choice |
| --- | --- |
| Surfacing behavior | Auto-copy newest email's code to clipboard + native notification |
| Detection strategy | Context-keyword required (a token is a code only if near a code keyword) |
| Settings control | New `AutoCopyCodes` toggle, default **ON** |
| Architecture | Isolated `internal/otp` package, wired into `app.go:handleEmail()` |

## Architecture

Approach A (chosen): a new, pure-Go `internal/otp` package exposing a single
detection function with no Wails dependency, called from the existing
`handleEmail` chokepoint. Clipboard via Wails' cross-platform
`runtime.ClipboardSetText`; notification via the existing `internal/notify`
package (already has `_darwin` / `_windows` / `_linux` builds). This mirrors the
existing pattern where `notify` and `settings` are separate internal packages
and keeps detection fully unit-testable.

Rejected:

- **B. Detect in `parser`, persist `OTPCode` on the model** — requires a DB
  migration and model change for UI display that is not in scope (YAGNI).
- **C. Frontend (JS) detection** — browser-style not native notifications, less
  reliable clipboard, and logic not covered by the Go test suite.

## Component Design

### 1. `internal/otp/otp.go`

```go
func Detect(email *models.Email) (code string, found bool)
```

- **Input text:** `email.Subject` + `email.PlainBody`. If `PlainBody` is empty,
  fall back to HTML-stripped text from `email.HTMLBody`.
- **Context keywords (case-insensitive):** `code`, `verification`, `verify`,
  `OTP`, `one-time`, `one time`, `2FA`, `two-factor`, `security`, `passcode`,
  `confirmation`, `confirm`, `authentication`, `authenticate`, `PIN`,
  `access code`, `login code`, `sign-in`, `sign in`.
- **Candidate token shapes:**
  - `4–8` consecutive digits.
  - Grouped digits: `\d{3}[- ]\d{3}` (normalized by stripping the separator).
  - `6–8` char alphanumeric containing at least one digit (uppercase-leaning;
    excludes plain dictionary words by requiring a digit).
- **Scoring:** locate keyword positions; a candidate qualifies only if it falls
  within a proximity window (~50 characters) of a keyword. The candidate with
  the smallest distance to a keyword wins; ties broken by earliest position.
- **False-positive exclusions:** bare 4-digit years in `1900–2099`;
  currency-adjacent tokens (immediately `$`-prefixed or `.dd`-suffixed).
- **No qualifying candidate → `found = false`.** The clipboard is never touched
  when detection is uncertain. Ambiguity means do not copy.
- Returned `code` is normalized (grouping separators stripped).

### 2. `app.go:handleEmail()`

After the existing `wailsruntime.EventsEmit(... "mailtap:email-received" ...)`:

```go
if settings.Load().AutoCopyCodes {
    if code, ok := otp.Detect(email); ok {
        if err := wailsruntime.ClipboardSetText(a.ctx, code); err != nil {
            log.Printf("Failed to copy OTP code to clipboard: %v", err)
        } else {
            go notify.Send("MailTap", fmt.Sprintf("Code %s copied to clipboard", code))
        }
    }
}
```

Best-effort and non-blocking; failures are logged, never fatal. The OTP
notification is governed solely by the new `AutoCopyCodes` toggle, independent
of the existing general `Notifications` setting (the two were chosen as one
combined behavior).

### 3. Settings

- Add `AutoCopyCodes bool` with json tag `autoCopyCodes` to the `Settings`
  struct in `internal/settings/settings.go`.
- `Default()` sets `AutoCopyCodes: true`.
- Backward compatibility: `Load()` starts from `Default()` then unmarshals the
  file over it, so existing `settings.json` files lacking the key retain the
  default value `true` (feature ON for existing users, matching the decision).
- `Settings.svelte`: add an "Auto-copy detected codes" toggle that mirrors the
  existing Notifications toggle, bound through the existing
  `GetSettings`/`SaveSettings` path (no new Wails bindings needed).

## Data Flow

1. SMTP server receives email → `parser.Parse()` → `models.Email`.
2. `App.handleEmail()` stores the email and emits `mailtap:email-received`.
3. If `AutoCopyCodes` is enabled, `otp.Detect(email)` runs.
4. On a confident detection: `runtime.ClipboardSetText` writes the code; a
   native notification "Code <code> copied to clipboard" is sent.
5. No detection or detection failure: nothing happens (silent, logged on error).

## Error Handling

- Detection never panics; malformed/empty bodies yield `found = false`.
- Clipboard write failure is logged; no notification is sent in that case.
- Notification send is fire-and-forget (`go notify.Send(...)`), as today.

## Testing

- **`internal/otp/otp_test.go`** — table-driven:
  - Positives: Google/GitHub/Stripe-style verification emails; 6-digit codes;
    grouped digits (`123-456`); alphanumeric codes with a digit.
  - Negatives: order-confirmation with an order number but no code keyword;
    price/date-only emails; alphanumeric tracking numbers without keywords;
    keyword present but no plausible token.
  - Normalization: grouped code returns separator-stripped value.
- **`internal/settings/settings_test.go`** (new or extended): assert
  `Default().AutoCopyCodes == true`, and that loading a legacy settings file
  without the `autoCopyCodes` key preserves `true`.

## Cross-Platform Notes

- `runtime.ClipboardSetText` is provided by Wails on macOS, Windows, and Linux.
- `internal/notify` already provides per-OS builds. The feature therefore works
  on Windows (the clipboard copy is fully cross-platform; the notification uses
  the existing platform-specific implementation).

## Out of Scope (YAGNI)

- Persisting the detected code on the `Email` model / in SQLite.
- A UI badge or in-detail "Copy code" button.
- Click-to-copy notification actions.
- Per-sender allow/deny lists or detection tuning UI.
