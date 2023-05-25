package executor

import (
	"email/dispatcher/executor/templates"
)

type SendEmailRequest struct {
	To      string
	Content templates.EmailContent
}
