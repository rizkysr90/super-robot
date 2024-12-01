package productservice

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"rizkysr90-pos/internal/payload"
	"rizkysr90-pos/internal/store"
	"rizkysr90-pos/pkg/errorHandler"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/signintech/gopdf"
)

const (
	barcodeWidth     = 1500 // 75mm in pixels
	barcodeHeight    = 350  // 35mm in pixels
	barcodesPerRow   = 3
	barcodeSpacingX  = 10.0
	barcodeSpacingY  = 50.0
	startXOffset     = 10.0
	startYOffset     = 10.0
	barcodeWidthMM   = 150.0
	barcodeHeightMM  = 35.0
	productNameSize  = 8
	productLabelSize = 6
	maxBarcodes      = 30
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

	// Generate and save barcode image
	tempFile, err := generateBarcode(product.ProductID)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile) // Clean up temporary file

	// Create PDF with barcodes
	pdfBytes, err := createPDF(product, tempFile)
	if err != nil {
		return nil, err
	}

	return &payload.GenerateBarcodeResponse{
		PDFBytes: pdfBytes,
	}, nil
}

func generateBarcode(productID string) (string, error) {
	// Generate barcode
	code, err := code128.Encode(productID)
	if err != nil {
		return "", fmt.Errorf("failed to generate barcode: %w", err)
	}

	// Scale barcode
	scaledBarcode, err := barcode.Scale(code, barcodeWidth, barcodeHeight)
	if err != nil {
		return "", fmt.Errorf("failed to scale barcode: %w", err)
	}

	// Convert to grayscale
	bounds := scaledBarcode.Bounds()
	grayscaleBarcode := image.NewGray(bounds)
	draw.Draw(grayscaleBarcode, bounds, scaledBarcode, bounds.Min, draw.Src)

	// Create unique filename with timestamp
	fileName := fmt.Sprintf("barcode_%s_%d.png", productID, time.Now().UnixNano())
	filePath := filepath.Join(os.TempDir(), fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create barcode file: %w", err)
	}
	defer file.Close()

	if err = png.Encode(file, grayscaleBarcode); err != nil {
		return "", fmt.Errorf("failed to encode barcode image: %w", err)
	}
	return filePath, nil
}

func createPDF(product *store.ProductData, barcodeImagePath string) ([]byte, error) {
	// Initialize PDF
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// Load font
	fontPath := filepath.Join("assets", "fonts", "Arial.ttf")
	if err := pdf.AddTTFFont("Arial", fontPath); err != nil {
		return nil, fmt.Errorf("failed to load font: %w", err)
	}

	// Add barcodes to PDF
	if err := addBarcodesToPDF(&pdf, product, barcodeImagePath); err != nil {
		return nil, err
	}

	// Write PDF to bytes
	var buf bytes.Buffer
	if _, err := pdf.WriteTo(&buf); err != nil {
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo("failed to generate PDF"),
			errorHandler.WithMessage(err.Error()),
		)
	}

	return buf.Bytes(), nil
}

func addBarcodesToPDF(pdf *gopdf.GoPdf, product *store.ProductData, imagePath string) error {
	for i := 0; i < maxBarcodes; i++ {
		xPos := startXOffset + (float64(i%barcodesPerRow) * (barcodeWidthMM + barcodeSpacingX))
		yPos := startYOffset + (float64(i/barcodesPerRow) * (barcodeHeightMM + barcodeSpacingY))

		// Add barcode image
		if err := pdf.Image(imagePath, xPos, yPos, &gopdf.Rect{W: barcodeWidthMM, H: barcodeHeightMM}); err != nil {
			return fmt.Errorf("failed to add barcode to PDF: %w", err)
		}

		// Add product name
		if err := addTextToPDF(pdf, product.ProductName, xPos, yPos+36, productNameSize); err != nil {
			return err
		}

		// Add price label
		priceLabel := fmt.Sprintf("%s - Rp%.2f", product.ProductID, product.Price)
		if err := addTextToPDF(pdf, priceLabel, xPos, yPos+43, productLabelSize); err != nil {
			return err
		}
	}
	return nil
}

func addTextToPDF(pdf *gopdf.GoPdf, text string, x, y float64, fontSize int) error {
	if err := pdf.SetFont("Arial", "", fontSize); err != nil {
		return fmt.Errorf("failed to set font: %w", err)
	}

	pdf.SetX(x)
	pdf.SetY(y)

	if err := pdf.Cell(nil, text); err != nil {
		return fmt.Errorf("failed to add text to PDF: %w", err)
	}

	return nil
}
