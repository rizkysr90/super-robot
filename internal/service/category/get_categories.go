package category

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/internal/store"
	"context"
	"math"

	"auth-service-rizkysr90-pos/pkg/errorHandler"
)

type reqGetCategories struct {
	*payload.ReqGetAllCategory
}

func (request *reqGetCategories) sanitize() {
	if request.PageNumber == 0 {
		request.PageNumber = 1
	}
	if request.PageSize == 0 {
		request.PageSize = 20
	}
}
func (c *Service) GetCategories(ctx context.Context,
	request *payload.ReqGetAllCategory) (*payload.ResGetAllCategory, error) {
	input := reqGetCategories{request}
	input.sanitize()

	categories, err := c.categoryStore.FindAllPagination(ctx,
		&store.Pagination{PageSize: input.PageSize, PageNumber: input.PageNumber})
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, errorHandler.NewNotFound(errorHandler.WithInfo("not found"))
	}
	result := &payload.ResGetAllCategory{
		Data: []payload.CategoryData{},
		Metadata: payload.Pagination{
			PageSize:      input.PageSize,
			PageNumber:    input.PageNumber,
			TotalPages:    int(math.Floor(float64(categories[0].Pagination.TotalElements) / float64(input.PageSize))),
			TotalElements: categories[0].Pagination.TotalElements,
		},
	}
	if result.Metadata.TotalPages == 0 {
		result.Metadata.TotalPages = 1
	}
	for _, element := range categories {
		data := payload.CategoryData{
			ID:           element.ID,
			CategoryName: element.CategoryName,
			CreatedAt:    element.CreatedAt,
			UpdatedAt:    element.UpdatedAt,
			DeletedAt:    element.DeletedAt.Time,
		}
		result.Data = append(result.Data, data)
	}
	return result, nil
}
