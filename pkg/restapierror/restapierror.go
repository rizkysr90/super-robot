package restapierror

import "context"

type RestAPIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func (err *RestAPIError) Error() string {
	return err.Message
}

type RestAPIErrorOption func(*RestAPIError)

func WithMessage(message string) RestAPIErrorOption {
	return func(err *RestAPIError) {
		err.Message = message
	}
}
func WithDetails(details interface{}) RestAPIErrorOption {
	return func(err *RestAPIError) {
		err.Details = details
	}
}
func NewBadRequest(ctx context.Context,
	opts ...RestAPIErrorOption) *RestAPIError {
	err := &RestAPIError{
		Code:    400,
		Message: "",
		Details: "",
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
