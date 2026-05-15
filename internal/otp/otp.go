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
