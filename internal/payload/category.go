package payload

import (
	"time"
)

type CategoryData struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
	ID           string    `json:"id"`
	CategoryName string    `json:"category_name"`
}
type ReqCreateCategory struct {
	CategoryName string `json:"category_name"`
}
type ResCreateCategory struct {
}
type ReqUpdateCategory struct {
	ID           string `json:"id"`
	CategoryName string `json:"category_name"`
}
type ResUpdateCategory struct {
}

type ReqGetCategoryByID struct {
	CategoryID string `json:"id"`
}
type ResGetCategoryByID struct {
	*CategoryData `json:"data"`
}
type ReqGetAllCategory struct {
	PageSize   int `json:"page_size"`
	PageNumber int `json:"page_number"`
}
type ResGetAllCategory struct {
	Data     []CategoryData `json:"data"`
	Metadata Pagination     `json:"metadata"`
}
type ReqDeleteCategory struct {
	ID string `json:"id"`
}
type ResDeleteCategory struct {
}
