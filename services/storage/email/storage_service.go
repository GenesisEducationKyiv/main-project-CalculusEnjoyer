package email

import (
	"fmt"
	"storage/domain"
	"storage/serror"

	"github.com/pkg/errors"
)

type Orchestrator interface {
	DeleteEmail(email domain.Email) error
	WriteEmail(email domain.Email) error
	GetAllRecords() ([]domain.Email, error)
}

type EmailRepository interface {
	Add(email domain.Email) error
	GetAll() (emails []domain.Email, err error)
	Delete(email domain.Email) error
	Exists(email domain.Email) (result bool, err error)
}

type fileEmailRepository struct {
	Orchestrator Orchestrator
}

func NewStorageRepository(orchestrator Orchestrator) EmailRepository {
	return &fileEmailRepository{Orchestrator: orchestrator}
}

func (r *fileEmailRepository) Add(email domain.Email) error {
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

func (r *fileEmailRepository) GetAll() ([]domain.Email, error) {
	return r.Orchestrator.GetAllRecords()
}

func (r *fileEmailRepository) Exists(email domain.Email) (bool, error) {
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

func (r *fileEmailRepository) Delete(email domain.Email) error {
	err := r.Orchestrator.DeleteEmail(email)
	return errors.Wrap(err, fmt.Sprintf("can not delete email %s", email.Value))
}
