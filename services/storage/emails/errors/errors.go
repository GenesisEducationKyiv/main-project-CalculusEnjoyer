package errors

import "fmt"

var (
	EmailAlreadyExist    = StorageError{Code: 0, Message: "email is already exists"}
	InternalStorageError = StorageError{Code: 1, Message: "unknown storage error"}
)

type StorageError struct {
	Message string
	Code    int
}

func (e StorageError) Error() string {
	return fmt.Sprintf("error: %s (code: %d)", e.Message, e.Code)
}
