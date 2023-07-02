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

type EmailRepository interface {
	Add(email messages.Email) (err error)
	GetAll() (emails []messages.Email, err error)
	Exists(email messages.Email) (result bool, err error)
}

type fileEmailRepository struct {
	Orchestrator Orchestrator
}

func NewStorageRepository(orchestrator Orchestrator) EmailRepository {
	return &fileEmailRepository{Orchestrator: orchestrator}
}

func (r *fileEmailRepository) Add(email messages.Email) error {
	isExist, err := r.Exists(email)
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

func (r *fileEmailRepository) GetAll() ([]messages.Email, error) {
	return r.Orchestrator.GetAllRecords()
}

func (r *fileEmailRepository) Exists(email messages.Email) (bool, error) {
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
