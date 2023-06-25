package ctrl

import (
	"api/models"
	"context"
)

type StorageOrchestrator interface {
	AddEmail(request models.AddEmailRequest, cnx context.Context) error
	GetAllEmails(cnx context.Context) ([]models.Email, error)
}
