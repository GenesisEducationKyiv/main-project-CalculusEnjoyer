package emails

import (
	"storage/emails/messages"
	"storage/serror"

	"github.com/pkg/errors"
)

type Orchestrator interface {
	WriteEmail(email messages.Email) error
	GetAllRecords() ([]messages.Email, error)
}

type StorageService interface {
	AddEmail(email messages.Email) (err error)
	GetAllEmails() (emails []messages.Email, err error)
	CheckIfEmailExists(email messages.Email) (result bool, err error)
}

type storageService struct {
	Orchestrator Orchestrator
}

func NewService(orchestrator Orchestrator) StorageService {
	return &storageService{Orchestrator: orchestrator}
}

func (r *storageService) AddEmail(email messages.Email) error {
	isExist, err := r.CheckIfEmailExists(email)
	if err != nil {
		return errors.Wrap(err, "can not check if email exists")
	}

	if isExist {
		return serror.ErrEmailAlreadyExists
	}

	err = r.Orchestrator.WriteEmail(email)
	if err != nil {
		return errors.Wrap(err, "can not write email")
	}

	return nil
}

func (r *storageService) GetAllEmails() ([]messages.Email, error) {
	return r.Orchestrator.GetAllRecords()
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
