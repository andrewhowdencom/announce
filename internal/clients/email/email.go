package email

import (
	"fmt"
	"net/smtp"
)

// Client is an interface for sending emails.
type Client interface {
	Send(to []string, author, subject, body string) error
}

// SMTPClient is a client for sending emails using SMTP.
type SMTPClient struct {
	addr string
	auth smtp.Auth
	from string
}

// NewClient creates a new SMTP client.
func NewClient(host string, port int, username, password, from string) Client {
	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%d", host, port)

	return &SMTPClient{
		addr: addr,
		auth: auth,
		from: from,
	}
}

// Send sends an email to the specified recipients.
func (c *SMTPClient) Send(to []string, author, subject, body string) error {
	var errs []error
	for _, recipient := range to {
		headers := map[string]string{
			"To":      recipient,
			"Subject": subject,
		}

		from := c.from
		if author != "" {
			from = author
			headers["Reply-To"] = author
		}
		headers["From"] = from

		msg := ""
		for k, v := range headers {
			msg += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		msg += "\r\n" + body

		err := smtp.SendMail(c.addr, c.auth, c.from, []string{recipient}, []byte(msg))
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to send email to %s: %w", recipient, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send email to some recipients: %v", errs)
	}

	return nil
}

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	SendFunc func(to []string, author, subject, body string) error
}

// NewMockClient returns a new mock client.
func NewMockClient() *MockClient {
	return &MockClient{}
}

// Send is the mock implementation of the Send method.
func (m *MockClient) Send(to []string, author, subject, body string) error {
	if m.SendFunc != nil {
		return m.SendFunc(to, author, subject, body)
	}
	return nil
}
