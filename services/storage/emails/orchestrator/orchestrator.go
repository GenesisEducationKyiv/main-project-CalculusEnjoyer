package orchestrator

import "storage/emails/messages"

type Orchestrator interface {
	WriteEmail(email messages.Email) error
	GetAllRecords() ([]messages.Email, error)
}
