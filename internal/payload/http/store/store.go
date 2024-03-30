package payload

type ReqCreateStore struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
	UserID  string `json:"user_id"`
}
