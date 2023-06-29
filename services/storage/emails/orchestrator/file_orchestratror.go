package orchestrator

import (
	"encoding/csv"
	"os"
	"storage/config"
	"storage/emails/messages"

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
	file, err := os.OpenFile(o.StoragePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
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

func (o *FileOrchestrator) GetAllRecords() ([]messages.Email, error) {
	cvsReader, err := o.ReadCSVData()
	if err != nil {
		return nil, errors.Wrap(err, "can not get reader")
	}

	allEmails, err := cvsReader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "can not get all records")
	}

	records := make([]messages.Email, len(allEmails))
	for i, emailRecord := range allEmails {
		records[i].Value = emailRecord[0]
	}

	return records, nil
}

func (o *FileOrchestrator) WriteEmail(email messages.Email) error {
	writer, err := o.WriteCsvData()
	if err != nil {
		return errors.Wrap(err, "can not write email")
	}

	defer writer.Flush()

	return writer.Write([]string{email.Value})
}
