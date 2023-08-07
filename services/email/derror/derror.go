package derror

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrSendEmails = status.Error(codes.Internal, "can not send email")
