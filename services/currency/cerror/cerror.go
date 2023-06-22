package cerror

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrRate   = status.Error(codes.InvalidArgument, "can not get crypto from third party")
	ErrDecode = status.Error(codes.Internal, "can not decode response from third party")
)

const ErrRateValue = -1
