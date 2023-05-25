package dispatcher

import (
	"email/dispatcher/executor"
)

type EmailService interface {
	SendEmail(req executor.SendEmailRequest) (err error)
}
