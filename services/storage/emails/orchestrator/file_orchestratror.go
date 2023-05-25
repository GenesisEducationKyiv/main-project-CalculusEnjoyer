package orchestrator

import (
	"encoding/csv"
	"github.com/joho/godotenv"
	"os"
)

var StoragePath string
var StorageFile *os.File

type FileOrchestrator struct{}

func init() {
	err := godotenv.Load()
	StoragePath = os.Getenv("STORAGE_PATH")

	if err != nil {
		print(err.Error())
	}
}

func (o *FileOrchestrator) OpenCSVFile() (*os.File, error) {

	file, err := os.OpenFile(StoragePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	StorageFile = file

	return StorageFile, nil
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
