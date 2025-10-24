package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/andrewhowdencom/ruf/internal/model"
)

// Client is an interface for sending emails.
type Client interface {
	Send(email *model.Email) error
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
func (c *SMTPClient) Send(email *model.Email) error {
	var errs []error
	for _, to := range email.To {
		msg := []byte(
			"To: " + to + "\r\n" +
				"Subject: " + email.Subject + "\r\n" +
				"\r\n" +
				email.Body + "\r\n",
		)

		err := smtp.SendMail(c.addr, c.auth, c.from, []string{to}, msg)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to send email to %s: %w", to, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send email to some recipients: %v", errs)
	}

	return nil
}

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	SendFunc func(email *model.Email) error
}

// NewMockClient returns a new mock client.
func NewMockClient() *MockClient {
	return &MockClient{}
}

// Send is the mock implementation of the Send method.
func (m *MockClient) Send(email *model.Email) error {
	if m.SendFunc != nil {
		return m.SendFunc(email)
	}
	return nil
}
