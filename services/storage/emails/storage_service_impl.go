package emails

import (
	"encoding/csv"
	storageErrors "storage/emails/errors"
	"storage/emails/messages"
	"storage/emails/orchestrator"

	"github.com/pkg/errors"
)

type storageService struct {
	Orchestrator orchestrator.FileOrchestrator
}

func NewService() StorageService {
	return &storageService{Orchestrator: *orchestrator.NewOrchestrator()}
}

func (r *storageService) AddEmail(email messages.Email) (err error) {
	isExist, err := r.CheckIfEmailExists(email)

	if err != nil {
		return errors.Wrap(err, "Error while checking if email already exist")
	}

	if !isExist && err == nil {
		var writer *csv.Writer
		writer, err = r.Orchestrator.WriteCsvData()

		if err != nil {
			return errors.Wrap(err, "Error while adding email")
		}

		err = writer.Write([]string{email.Value})

		if err != nil {
			return errors.Wrap(err, "Error while adding email")
		}

		writer.Flush()
	}

	if isExist {
		err = storageErrors.EmailAlreadyExist()
	}

	return err
}

func (r *storageService) GetAllEmails() (emails []messages.Email, err error) {
	allData, err := r.Orchestrator.GetAllRecords()

	if err != nil {
		return nil, errors.Wrap(err, "Error while getting all emails")
	}

	for i := range allData {
		emails = append(emails, messages.Email{Value: allData[i][0]})
	}

	return emails, err
}

func (r *storageService) CheckIfEmailExists(email messages.Email) (result bool, err error) {
	allData, err := r.Orchestrator.GetAllRecords()

	if err != nil {
		return false, errors.Wrap(err, "Error while getting all emails")
	}

	for i := range allData {
		if (allData[i][0]) == email.Value {
			return true, err
		}
	}

	return false, err
}
