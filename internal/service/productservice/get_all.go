package productservice

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/internal/store"
	"auth-service-rizkysr90-pos/pkg/errorHandler"
	"context"
)


func (s *Service) GetAllProducts(ctx context.Context, 
	request *payload.ReqGetAllProducts) (*payload.ResGetAllProducts, error) {
	// Validate pagination parameters
	if request.PageSize <= 0 {
		request.PageSize = 10 // Default page size
	}
	if request.PageNumber <= 0 {
		request.PageNumber = 1 // Default page number
	}
	// Get products from the store
	products, totalCount, err := s.productStore.GetAll(ctx, &store.FilterProduct{
		Limit:      request.PageSize,
		Offset:     (request.PageNumber - 1) * request.PageSize,
		CategoryID: request.CategoryID,
	})
	if err != nil {
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo("failed to get products"),
			errorHandler.WithMessage(err.Error()),
		)
	}
	// Prepare response
	response := &payload.ResGetAllProducts{
		Data: make([]payload.ProductData, len(products)),
		Metadata: payload.Pagination{
			PageSize:      request.PageSize,
			PageNumber:    request.PageNumber,
			TotalPages:    (totalCount + request.PageSize - 1) / request.PageSize,
			TotalElements: totalCount,
		},
	}
	for i, p := range products {
		response.Data[i] = payload.ProductData{
			ProductID:     p.ProductID,
			ProductName:   p.ProductName,
			Price:         p.Price,
			BasePrice:     p.BasePrice,
			StockQuantity: p.StockQuantity,
			CategoryID:    p.Category.ID,
			CategoryName:  p.Category.CategoryName,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		}
	}
	return response, nil
}