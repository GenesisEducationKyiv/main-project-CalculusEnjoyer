package ctrl

import (
	"api/models"
	"context"
)

type EmailExecutor interface {
	SendEmail(request models.SendEmailsRequest, cnx context.Context) error
}
