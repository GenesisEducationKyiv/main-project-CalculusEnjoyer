package emails

import (
	storageErrors "storage/emails/error"
	"storage/emails/messages"
	"storage/emails/orchestrator"
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
		return storageErrors.InternalStorageError()
	}

	if !isExist && err == nil {
		err = r.Orchestrator.WriteEmail(email)
		if err != nil {
			return storageErrors.InternalStorageError()
		}
	}

	if isExist {
		err = storageErrors.EmailAlreadyExist()
	}

	return err
}

func (r *storageService) GetAllEmails() (emails []messages.Email, err error) {
	allData, err := r.Orchestrator.GetAllRecords()
	if err != nil {
		return nil, storageErrors.InternalStorageError()
	}

	for i := range allData {
		emails = append(emails, messages.Email{Value: allData[i][0]})
	}

	return emails, err
}

func (r *storageService) CheckIfEmailExists(email messages.Email) (result bool, err error) {
	allData, err := r.Orchestrator.GetAllRecords()
	if err != nil {
		return false, storageErrors.InternalStorageError()
	}

	for i := range allData {
		if (allData[i][0]) == email.Value {
			return true, err
		}
	}

	return false, err
}
