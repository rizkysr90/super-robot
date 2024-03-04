package restapierror

import "context"

type RestAPIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"more_details"`
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
func NewInternalServer(ctx context.Context,
	opts ...RestAPIErrorOption) *RestAPIError {

	err := &RestAPIError{
		Code:    500,
		Message: "Internal Server Error",
		Details: "",
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
func NewMultipleFieldsValidation(ctx context.Context, errors []RestAPIError, opts ...RestAPIErrorOption) *RestAPIError {
	err := &RestAPIError{
		Code:    400,
		Message: "Bad Request - Error Multiple Fields Validaition",
		Details: errors,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
