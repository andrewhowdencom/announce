package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/andrewhowdencom/ruf/internal/config"
)

// EmailClient is an interface for sending emails.
type EmailClient interface {
	Send(to []string, subject, body string) error
}

// SMTPClient is an implementation of EmailClient that sends emails using SMTP.
type SMTPClient struct {
	cfg config.Email
}

// NewSMTPClient returns a new SMTPClient.
func NewSMTPClient(cfg config.Email) *SMTPClient {
	return &SMTPClient{cfg: cfg}
}

// Send sends an email using SMTP.
func (c *SMTPClient) Send(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)

	headers := make(map[string]string)
	headers["From"] = c.cfg.From
	headers["To"] = strings.Join(to, ", ")
	headers["Subject"] = subject

	var msg string
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
	return smtp.SendMail(addr, auth, c.cfg.From, to, []byte(msg))
}
