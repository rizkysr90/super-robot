package store

import (
	"context"
	"database/sql"
	"time"
)

type EmployeeData struct {
	ID           string
	Name         string
	Contact      string
	Username     string
	Password     string
	StoreID      string
	RefreshToken string
	Role         int
	CreatedAt    time.Time
	DeletedAt    sql.NullTime
}
type SetResponsGetEmployees struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Contact   string    `json:"contact"`
	Username  string    `json:"username"`
	StoreID   string    `json:"store_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type EmployeeFilter struct {
	Username string
	StoreID  string
}

// func (s *StoreData) ToPayloadResponse(data *StoreData) *SetResponse {
// 	return &SetResponse{
// 		ID:        s.ID,
// 		Name:      s.Name,
// 		Address:   s.Address,
// 		Contact:   s.Contact,
// 		CreatedAt: s.CreatedAt,
// 		DeletedAt: s.DeletedAt.Time,
// 	}
// }

type EmployeeStore interface {
	Insert(ctx context.Context, data *EmployeeData) error
	Finder(ctx context.Context, filter *EmployeeFilter, operationID string) ([]EmployeeData, error)
	FindOne(ctx context.Context, filterBy *EmployeeFilter, operationID string) (*EmployeeData, error)
	Update(ctx context.Context,
		updatedData *EmployeeData,
		filter *EmployeeFilter,
		staging string) error
}
