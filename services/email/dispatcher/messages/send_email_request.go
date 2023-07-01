package messages

type SendEmailRequest struct {
	To      string
	Content EmailContent
}
