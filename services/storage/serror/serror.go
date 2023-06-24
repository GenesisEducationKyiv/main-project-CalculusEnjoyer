package serror

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternalStorageError = status.Error(codes.Internal, "internal storage error")
	ErrEmailAlreadyExists   = status.Error(codes.AlreadyExists, "email already exist")
	ErrStorage              = status.Error(codes.Internal, "storage is not working properlly")
)
