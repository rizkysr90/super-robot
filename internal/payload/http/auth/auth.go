package payload

type ReqCreateAccount struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
type ReqLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ResLoginUser struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
