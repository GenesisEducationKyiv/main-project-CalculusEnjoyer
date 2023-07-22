package domain

type SendEmailRequest struct {
	To      string
	Content EmailContent
}
