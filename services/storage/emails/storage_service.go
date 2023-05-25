package emails

import (
	"storage/emails/messages"
)

type StorageService interface {
	AddEmail(email messages.Email) (err error)
	GetAllEmails() (emails []messages.Email, err error)
	CheckIfEmailExists(email messages.Email) (result bool, err error)
}
