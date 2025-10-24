package worker

import (
	"fmt"

	"github.com/andrewhowdencom/ruf/internal/email"
	"github.com/andrewhowdencom/ruf/internal/model"
)

// EmailWorker is a worker that sends emails.
type EmailWorker struct {
	emailClient email.EmailClient
}

// NewEmailWorker returns a new EmailWorker.
func NewEmailWorker(emailClient email.EmailClient) *EmailWorker {
	return &EmailWorker{emailClient: emailClient}
}

// Process processes a call and sends an email if the call has an email defined.
func (w *EmailWorker) Process(call model.Call) error {
	if call.Email == nil {
		return nil
	}

	fmt.Printf("Sending email for call %s to %v with subject '%s'...", call.ID, call.Email.To, call.Email.Subject)

	err := w.emailClient.Send(call.Email.To, call.Email.Subject, call.Email.Body)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
