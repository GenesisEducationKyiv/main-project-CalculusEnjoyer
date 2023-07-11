package aerror

import "errors"

var (
	ErrGRPC                   = errors.New("gRPC interceptor dramatically collapsed")
	ErrUnsupportedCurrency    = errors.New("unsupported currency")
	ErrRequest                = errors.New("can not decode response")
	ErrInvalidEmail           = errors.New("invalid email")
	ErrFailedToSend           = errors.New("failed to send email")
	ErrFailedToEncodeResponse = errors.New("failed to encode response")
)
