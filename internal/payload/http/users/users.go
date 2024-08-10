package users

type ReqCreateUsers struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
type ResCreateUsers struct {

}

type ResLoginUsers struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
}