package payload

import (
	"auth-service-rizkysr90-pos/internal/store"
)

type ReqCreateStore struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
	UserID  string `json:"user_id"`
}
type ReqGetAllStore struct {
	UserID     string // get from context
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
}

type ResGetAllStore struct {
	Data       []store.SetResponse
	Pagination store.Pagination
}
