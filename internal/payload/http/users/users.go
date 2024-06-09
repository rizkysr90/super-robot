package users

type ReqCreateUsers struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Phone           string `json:"phone"`
	Role            int    `json:"role"`
}
type ReqLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ResLoginUser struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
type ReqRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}
type ResRefreshToken struct {
	// it will attach to the cookies, so response body remain empty
	AccessToken string
}
