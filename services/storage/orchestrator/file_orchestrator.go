package orchestrator

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"storage/config"
	"storage/domain"

	"github.com/pkg/errors"
)

type FileOrchestrator struct {
	StoragePath string
}

func NewFileOrchestrator(config config.Config) *FileOrchestrator {
	orchestrator := FileOrchestrator{StoragePath: config.EmailStoragePath}

	return &orchestrator
}

func (o *FileOrchestrator) OpenCSVFile() (*os.File, error) {
	file, err := os.OpenFile(o.StoragePath, os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (o *FileOrchestrator) ReadCSVData() (*csv.Reader, error) {
	file, err := o.OpenCSVFile()
	if err != nil {
		return nil, errors.Wrap(err, "can not open file")
	}

	return csv.NewReader(file), nil
}

func (o *FileOrchestrator) WriteCsvData() (*csv.Writer, error) {
	file, err := o.OpenCSVFile()
	if err != nil {
		return nil, errors.Wrap(err, "can not open file")
	}

	return csv.NewWriter(file), nil
}

func (o *FileOrchestrator) GetAllRecords() ([]domain.Email, error) {
	cvsReader, err := o.ReadCSVData()
	if err != nil {
		return nil, errors.Wrap(err, "can not get reader")
	}

	allEmails, err := cvsReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "can not get all records")
	}

	records := make([]domain.Email, len(allEmails))
	for i, emailRecord := range allEmails {
		records[i].Value = emailRecord[0]
	}

	return records, nil
}

func (o *FileOrchestrator) WriteEmail(email domain.Email) error {
	writer, err := o.WriteCsvData()
	if err != nil {
		return errors.Wrap(err, "can not write email")
	}

	defer writer.Flush()

	return writer.Write([]string{email.Value})
}

func (o *FileOrchestrator) DeleteEmail(target domain.Email) (err error) {
	emails, err := o.GetAllRecords()
	if err != nil {
		return errors.Wrap(err, "can not get emails while deleting email")
	}

	o.clearTestFile()

	emailsNoTarget, err := deleteEmailFromSlice(&emails, target)
	if err != nil {
		return errors.Wrap(err, "can not remove email")
	}

	for _, email := range *emailsNoTarget {
		writeError := o.WriteEmail(email)
		if writeError != nil {
			if err == nil {
				err = writeError
			}
			log.Printf("error while coping %s", email.Value)
			err = errors.Wrap(err, fmt.Sprintf("error while coping %s", email.Value))
		}
	}

	return err
}

func (o *FileOrchestrator) clearTestFile() {
	if err := os.Truncate(o.StoragePath, 0); err != nil {
		log.Fatalf(errors.Wrap(err, "Can not load test config").Error())
	}
}

func deleteEmailFromSlice(emails *[]domain.Email, target domain.Email) (*[]domain.Email, error) {
	index := -1
	for i, email := range *emails {
		if email.Value == target.Value {
			index = i
			break
		}
	}

	if index == -1 {
		return nil, errors.New(fmt.Sprintf("in storage there is no %s", target))
	}

	result := append((*emails)[:index], (*emails)[index+1:]...)
	return &result, nil
}
