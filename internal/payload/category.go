package payload

import (
	"time"
)

type CategoryData struct {
	Id           string    `json:"id"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}
type ReqCreateCategory struct {
	CategoryName        string `json:"category_name"`
}
type ResCreateCategory struct {

}
type ReqUpdateCategory struct {
	Id string `json:"id"`
	CategoryName        string `json:"category_name"`
}
type ResUpdateCategory struct {
}

type ReqGetCategoryById struct {
	Id string `json:"id"`
}
type ResGetCategoryById struct {
	*CategoryData `json:"data"`
}
type ReqGetAllCategory struct {
	PageSize int `json:"page_size"`
	PageNumber int `json:"page_number"`
}
type ResGetAllCategory struct {
	Data     []CategoryData      `json:"data"`
	Metadata Pagination `json:"metadata"`
}