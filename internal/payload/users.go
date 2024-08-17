package payload

type ReqCreateUsers struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
type ResCreateUsers struct {

}

type ReqLoginUsers struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
}
type ResLoginUsers struct {
	AccessToken string `json:"access_token"`
}