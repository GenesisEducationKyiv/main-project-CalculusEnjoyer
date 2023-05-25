package emails

import (
	"encoding/csv"
	storageErrors "storage/emails/errors"
	"storage/emails/messages"
	"storage/emails/orchestrator"
)

type storageService struct {
	Orchestrator orchestrator.FileOrchestrator
}

func NewService() StorageService {
	return &storageService{Orchestrator: orchestrator.FileOrchestrator{}}
}

func (r *storageService) AddEmail(email messages.Email) (err error) {
	isExist, err := r.CheckIfEmailExists(email)

	if !isExist && err == nil {
		var writer *csv.Writer
		writer, err = r.Orchestrator.WriteCsvData()
		err = writer.Write([]string{email.Value})
		writer.Flush()
	}

	if isExist {
		err = storageErrors.EmailAlreadyExist
	}

	return err
}

func (r *storageService) GetAllEmails() (emails []messages.Email, err error) {
	cvsReader, err := r.Orchestrator.ReadCSVData()
	allData, err := (*cvsReader).ReadAll()

	for i := range allData {
		emails = append(emails, messages.Email{Value: allData[i][0]})
	}

	return emails, err
}

func (r *storageService) CheckIfEmailExists(email messages.Email) (result bool, err error) {
	cvsReader, err := r.Orchestrator.ReadCSVData()
	allData, err := (*cvsReader).ReadAll()

	if err != nil {
		return false, err
	}

	for i := range allData {
		if (allData[i][0]) == email.Value {
			return true, err
		}
	}

	return false, err
}
