package emails

import (
	"storage/emails/messages"
	"storage/emails/orchestrator"
	"storage/serror"

	"github.com/pkg/errors"
)

type storageService struct {
	Orchestrator orchestrator.Orchestrator
}

func NewService(orchestrator orchestrator.Orchestrator) StorageService {
	return &storageService{Orchestrator: orchestrator}
}

func (r *storageService) AddEmail(email messages.Email) error {
	isExist, err := r.CheckIfEmailExists(email)
	if err != nil {
		return errors.Wrap(err, "can not check if email exists")
	}

	if !isExist {
		err = r.Orchestrator.WriteEmail(email)
		if err != nil {
			return errors.Wrap(err, "can not write email")
		}
	}

	if isExist {
		return serror.ErrEmailAlreadyExists
	}

	return nil
}

func (r *storageService) GetAllEmails() (emails []messages.Email, err error) {
	allEmails, err := r.Orchestrator.GetAllRecords()
	if err != nil {
		return nil, errors.Wrap(err, "can not get all emails")
	}

	emails = append(emails, allEmails...)

	return emails, nil
}

func (r *storageService) CheckIfEmailExists(email messages.Email) (bool, error) {
	allData, err := r.Orchestrator.GetAllRecords()
	if err != nil {
		return false, errors.Wrap(err, "can not get all records")
	}

	for i := range allData {
		if (allData[i].Value) == email.Value {
			return true, nil
		}
	}

	return false, nil
}
