package constant

const (
	ErrNotFound       = "no data found"
	ErrInternalServer = "internal server error"
	ErrInvalidAPIKey  = "Unauthorized: Invalid API key"
	ErrInvalidFormat  = "Invalid Format"
)

type UserType string

const (
	UserTypeGoogle UserType = "google"
)
