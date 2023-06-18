package orchestrator

import (
	"encoding/csv"
	"os"
	StorageError "storage/emails/error"
	"storage/emails/messages"

	"github.com/joho/godotenv"
)

type FileOrchestrator struct {
	StoragePath string
	StorageFile *os.File
}

func NewOrchestrator() *FileOrchestrator {
	err := godotenv.Load()
	orchestrator := FileOrchestrator{StoragePath: os.Getenv("STORAGE_PATH")}

	if err != nil {
		print(err.Error())
	}

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
		return nil, err
	}

	return csv.NewReader(file), nil
}

func (o *FileOrchestrator) WriteCsvData() (*csv.Writer, error) {
	file, err := o.OpenCSVFile()
	if err != nil {
		return nil, err
	}

	return csv.NewWriter(file), nil
}

func (o *FileOrchestrator) GetAllRecords() (records [][]string, err error) {
	cvsReader, err := o.ReadCSVData()
	if err != nil {
		return nil, StorageError.InternalStorageError()
	}

	allData, err := (*cvsReader).ReadAll()
	if err != nil {
		return nil, StorageError.InternalStorageError()
	}

	return allData, nil
}

func (o *FileOrchestrator) WriteEmail(email messages.Email) (err error) {
	var writer *csv.Writer
	writer, err = o.WriteCsvData()

	defer writer.Flush()
	if err != nil {
		return StorageError.InternalStorageError()
	}

	return writer.Write([]string{email.Value})
}
