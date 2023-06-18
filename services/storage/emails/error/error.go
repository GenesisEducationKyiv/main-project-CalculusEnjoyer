package error

import "fmt"

func EmailAlreadyExist() StorageError {
	return StorageError{Code: 0, Message: "email is already exists"}
}

func InternalStorageError() StorageError {
	return StorageError{Code: 1, Message: "unknown storage error"}
}

type StorageError struct {
	Message string
	Code    int
}

func (e StorageError) Error() string {
	return fmt.Sprintf("error: %s (code: %d)", e.Message, e.Code)
}
