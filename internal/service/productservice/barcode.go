package productservice

import (
	"context"
	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/pkg/errorHandler"
)

// GenerateBarcodePDF generates a PDF containing barcodes for the given product.
func (s *Service) GenerateBarcodePDF(ctx context.Context,
	req *payload.GenerateBarcodeRequest) (*payload.GenerateBarcodeResponse, error) {
	// Fetch product details
	product, err := s.productStore.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo("failed to fetch product"),
			errorHandler.WithMessage(err.Error()),
		)
	}
	// Create PDF with barcodes
	pdfBytes, err := s.documentGenerator.LabelPricing(product)
	if err != nil {
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo("failed to generate barcode product"),
			errorHandler.WithMessage(err.Error()),
		)
	}
	return &payload.GenerateBarcodeResponse{
		PDFBytes: pdfBytes,
	}, nil
}
