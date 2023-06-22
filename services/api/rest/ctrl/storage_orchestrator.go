package ctrl

import (
	"api/models"
)

type StorageOrchestrator interface {
	AddEmail(request models.AddEmailRequest) error
	GetAllEmails() ([]models.Email, error)
}
