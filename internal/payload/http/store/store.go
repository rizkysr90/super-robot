package payload

import (
	"auth-service-rizkysr90-pos/internal/store"
)

type ReqCreateStore struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	Contact    string `json:"contact"`
	EmployeeID string `json:"user_id"`
}
type ReqUpdateStore struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	Contact    string `json:"contact"`
	EmployeeID string // get from jwt session
	StoreID    string `json:"store_id"`
}
type ResCreateStore struct {
	StoreID   string `json:"store_id"`
	CreatedAt string `json:"created_at"`
}
type ReqGetAllStore struct {
	EmployeeID string // get from context
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
}
type ReqDeleteStore struct {
	EmployeeID string // get from context
	StoreID    string `json:"store_id"`
}

type ResGetAllStore struct {
	Data       []store.SetResponse
	Pagination store.Pagination
}
