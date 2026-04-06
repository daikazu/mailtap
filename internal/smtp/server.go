package smtp

import (
	"bytes"
	"io"
	"log"
	"net"
	"time"

	gosmtp "github.com/emersion/go-smtp"

	"mailtap/internal/models"
	"mailtap/internal/parser"
)

// EmailHandler is a callback function invoked when an email is received.
type EmailHandler func(email *models.Email)

// Server wraps the go-smtp server and manages the listener lifecycle.
type Server struct {
	server   *gosmtp.Server
	handler  EmailHandler
	listener net.Listener
}

// NewServer creates a new SMTP server bound to the given address.
// The handler is called for each successfully received email.
func NewServer(addr string, handler EmailHandler) (*Server, error) {
	s := &Server{handler: handler}

	be := &backend{handler: handler}
	s.server = gosmtp.NewServer(be)
	s.server.Addr = addr
	s.server.Domain = "localhost"
	s.server.ReadTimeout = 30 * time.Second
	s.server.WriteTimeout = 30 * time.Second
	s.server.MaxMessageBytes = 25 * 1024 * 1024
	s.server.MaxRecipients = 100
	s.server.AllowInsecureAuth = true

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	s.listener = ln

	return s, nil
}

// Start begins accepting SMTP connections. It blocks until the server is stopped.
func (s *Server) Start() error {
	log.Printf("SMTP server listening on %s", s.listener.Addr().String())
	return s.server.Serve(s.listener)
}

// Stop gracefully shuts down the SMTP server.
func (s *Server) Stop() error {
	return s.server.Close()
}

// Addr returns the address the server is listening on.
func (s *Server) Addr() string {
	return s.listener.Addr().String()
}

// backend implements the gosmtp.Backend interface.
type backend struct {
	handler EmailHandler
}

func (b *backend) NewSession(_ *gosmtp.Conn) (gosmtp.Session, error) {
	return &session{handler: b.handler}, nil
}

// session implements the gosmtp.Session interface.
type session struct {
	handler EmailHandler
	from    string
	to      []string
}

func (s *session) Mail(from string, _ *gosmtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *session) Rcpt(to string, _ *gosmtp.RcptOptions) error {
	s.to = append(s.to, to)
	return nil
}

func (s *session) Data(r io.Reader) error {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	email, err := parser.Parse(buf.Bytes())
	if err != nil {
		log.Printf("Warning: failed to parse email: %v", err)
		return nil
	}

	if email.From == "" {
		email.From = s.from
	}
	if len(email.To) == 0 {
		email.To = s.to
	}

	s.handler(email)
	return nil
}

func (s *session) Reset() {
	s.from = ""
	s.to = nil
}

func (s *session) Logout() error {
	return nil
}
