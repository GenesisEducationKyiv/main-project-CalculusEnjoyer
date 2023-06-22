package models

import "email/dispatcher/executor/templates"

type SendEmailsRequest struct {
	Interceptor Email
	Template    templates.EmailContent
}
