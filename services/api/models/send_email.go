package models

import (
	"email/dispatcher/messages"
)

type SendEmailsRequest struct {
	Interceptor Email
	Template    messages.EmailContent
}
