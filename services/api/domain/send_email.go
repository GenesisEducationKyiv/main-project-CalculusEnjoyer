package domain

type SendEmailsRequest struct {
	Interceptor Email
	Template    EmailContent
}

type EmailContent struct {
	Body    string
	Subject string
}
