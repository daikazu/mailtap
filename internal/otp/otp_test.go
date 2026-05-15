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
