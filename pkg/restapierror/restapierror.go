package restapierror

import "github.com/rizkysr90/go-boilerplate/internal/constant"

type RestAPIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"more_details"`
}

func (err *RestAPIError) Error() string {
	return err.Message
}

type Option func(*RestAPIError)

func WithMessage(message string) Option {
	return func(err *RestAPIError) {
		err.Message = message
	}
}
func WithDetails(details interface{}) Option {
	return func(err *RestAPIError) {
		err.Details = details
	}
}
func NewBadRequest(
	opts ...Option) *RestAPIError {
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
func NewInternalServer(
	opts ...Option) *RestAPIError {
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
func NewMultipleFieldsValidation(errors []RestAPIError, opts ...Option) *RestAPIError {
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
func NewUnauthorized(opts ...Option) *RestAPIError {
	err := &RestAPIError{
		Code:    401,
		Message: constant.ErrInvalidAPIKey,
		Details: "",
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
func NewNotFound(opts ...Option) *RestAPIError {
	err := &RestAPIError{
		Code:    404,
		Message: constant.ErrNotFound,
		Details: "",
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
