package payload

type ReqCreateEmployee struct {
	Name            string `json:"name"`
	Contact         string `json:"contact"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	StoreID         string `json:"store_id"`
	Role            int    `json:"role"`
}
type ReqLoginEmployee struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type ResLoginEmployee struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	Role         int    `json:"role"`
}
