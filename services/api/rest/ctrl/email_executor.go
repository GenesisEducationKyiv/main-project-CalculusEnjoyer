package ctrl

import (
	"api/models"
)

type EmailExecutor interface {
	SendEmail(request models.SendEmailsRequest) error
}
