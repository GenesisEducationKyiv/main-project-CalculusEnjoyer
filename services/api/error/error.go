package error

import "fmt"

func RateProviderError() APIerror {
	return APIerror{Code: 2, Message: "can not get rate from rate provider service"}
}

func InvalidEmail() APIerror {
	return APIerror{Code: 3, Message: "invalid email"}
}

func InvalidRequest() APIerror {
	return APIerror{Code: 4, Message: "invalid request"}
}

type APIerror struct {
	Message string
	Code    int
}

func (e APIerror) Error() string {
	return fmt.Sprintf("error: %s (code: %d)", e.Message, e.Code)
}
