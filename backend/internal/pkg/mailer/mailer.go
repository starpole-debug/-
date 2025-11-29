package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

// Client sends plain-text emails via SMTP.
type Client interface {
	Send(to, subject, body string) error
}

// SMTPClient implements Client using STARTTLS where available.
type SMTPClient struct {
	host         string
	addr         string
	auth         smtp.Auth
	fromHeader   string
	fromEnvelope string
}

// NewSMTPClient builds a mailer; returns nil if host is empty.
func NewSMTPClient(host string, port int, username, password, from string) Client {
	host = strings.TrimSpace(host)
	if host == "" {
		return nil
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	var auth smtp.Auth
	if username != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}
	return &SMTPClient{
		host:         host,
		addr:         addr,
		auth:         auth,
		fromHeader:   from,
		fromEnvelope: envelopeAddress(from, username),
	}
}

func (c *SMTPClient) Send(to, subject, body string) error {
	if c == nil {
		return fmt.Errorf("mailer not configured")
	}
	msg := buildMessage(c.fromHeader, to, subject, body)

	conn, err := smtp.Dial(c.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// STARTTLS if supported
	if ok, _ := conn.Extension("STARTTLS"); ok {
		tlsCfg := &tls.Config{
			ServerName: c.host,
		}
		if err := conn.StartTLS(tlsCfg); err != nil {
			return err
		}
	}

	if c.auth != nil {
		if ok, _ := conn.Extension("AUTH"); ok {
			if err := conn.Auth(c.auth); err != nil {
				return err
			}
		}
	}

	if err := conn.Mail(c.fromEnvelope); err != nil {
		return err
	}
	if err := conn.Rcpt(to); err != nil {
		return err
	}

	wc, err := conn.Data()
	if err != nil {
		return err
	}
	if _, err := wc.Write([]byte(msg)); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return conn.Quit()
}

func buildMessage(from, to, subject, body string) string {
	headers := []string{
		"From: " + from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-version: 1.0;",
		"Content-Type: text/plain; charset=\"UTF-8\";",
		"",
	}
	return strings.Join(headers, "\r\n") + "\r\n" + body
}

// envelopeAddress extracts the raw email address for SMTP MAIL FROM.
func envelopeAddress(from, fallback string) string {
	from = strings.TrimSpace(from)
	if from == "" {
		return fallback
	}
	if strings.Contains(from, "<") && strings.Contains(from, ">") {
		start := strings.Index(from, "<")
		end := strings.Index(from, ">")
		if start >= 0 && end > start {
			return strings.TrimSpace(from[start+1 : end])
		}
	}
	if strings.Contains(from, " ") {
		// Likely a display name without angle brackets; fall back to username.
		if fallback != "" {
			return fallback
		}
	}
	return from
}
