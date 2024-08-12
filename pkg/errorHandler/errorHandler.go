package errorHandler

const (
	ErrInvalidAPIKey = "API KEY IS INVALID"
	ErrNotFound      = "RESOURCES NOT FOUND"
)

type HttpError struct {
	Code    int         `json:"code"`
	Info    string      `json:"info"`
	Message interface{} `json:"message"`
}

func (err *HttpError) Error() string {
	return err.Info
}

type Option func(*HttpError)

func WithMessage(message interface{}) Option {
	return func(err *HttpError) {
		err.Message = message
	}
}
func WithInfo(info string) Option {
	return func(err *HttpError) {
		err.Info = info
	}
}
func NewBadRequest(
	opts ...Option) *HttpError {
	err := &HttpError{
		Code:    400,
		Info:    "Bad Request",
		Message: "",
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
func NewInternalServer(
	opts ...Option) *HttpError {
	err := &HttpError{
		Code:    500,
		Info:    "Internal Server Error",
		Message: "",
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
func NewMultipleFieldsValidation(errors []HttpError, opts ...Option) *HttpError {
	err := &HttpError{
		Code:    400,
		Info:    "Bad Request - Error Multiple Fields Validaition",
		Message: errors,
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
func NewUnauthorized(opts ...Option) *HttpError {
	err := &HttpError{
		Code:    401,
		Info:    ErrInvalidAPIKey,
		Message: "",
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
func NewNotFound(opts ...Option) *HttpError {
	err := &HttpError{
		Code:    404,
		Info:    ErrNotFound,
		Message: "",
	}
	for _, opt := range opts {
		opt(err)
	}
	return err
}
