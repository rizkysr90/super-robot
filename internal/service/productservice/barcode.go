package productservice

import (
	"auth-service-rizkysr90-pos/internal/payload"
	"auth-service-rizkysr90-pos/pkg/errorHandler"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/signintech/gopdf"
)

func (s *Service) GenerateBarcodePDF(ctx context.Context, req *payload.GenerateBarcodeRequest) (*payload.GenerateBarcodeResponse, error) {
	// Fetch product details
	product, err := s.productStore.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo("failed to fetch product"),
			errorHandler.WithMessage(err.Error()),
		)
	}

	// Generate barcode from product ID (assuming ProductID is a string)
	code, err := code128.Encode(product.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate barcode: %v", err)
	}

	// Scale the barcode to the required size (75mm x 35mm)
	barcodeWidth := 1500 // 75 mm in pixels
	barcodeHeight := 350 // 35 mm in pixels
	scaledBarcode, err := barcode.Scale(code, barcodeWidth, barcodeHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to scale barcode: %v", err)
	}

	// Convert the barcode image to grayscale
	bounds := scaledBarcode.Bounds()
	grayscaleBarcode := image.NewGray(bounds)
	draw.Draw(grayscaleBarcode, bounds, scaledBarcode, bounds.Min, draw.Src)

	// Save the grayscale barcode image to a temporary file
	tempFile := "barcode.png"
	file, err := os.Create(tempFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create barcode file: %v", err)
	}
	defer file.Close()

	err = png.Encode(file, grayscaleBarcode)
	if err != nil {
		return nil, fmt.Errorf("failed to encode barcode image: %v", err)
	}

	// Create a new PDF
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4}) // A4 page size
	pdf.AddPage()

	// Set up barcode placement on the A4 page
	xOffset := 10.0 // starting X position in mm
	yOffset := 10.0 // starting Y position in mm
	barcodesPerRow := 3
	barcodeSpacingX := 10.0 // space between barcodes horizontally
	barcodeSpacingY := 50.0 // space between barcodes vertically

	// Load a font for the labels
	fontPath := filepath.Join("assets", "fonts", "Arial.ttf")
	err = pdf.AddTTFFont("Arial", fontPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %v", err)
	}

	// Load the barcode image into the PDF
	for i := 0; i < 30; i++ { // adjust this loop based on the number of barcodes you want
		xPos := xOffset + (float64(i%barcodesPerRow) * (150.0 + barcodeSpacingX)) // adjust X position for each barcode
		yPos := yOffset + (float64(i/barcodesPerRow) * (35.0 + barcodeSpacingY)) // adjust Y position for each row

		err = pdf.Image(tempFile, xPos, yPos, &gopdf.Rect{W: 150, H: 35}) // 76 mm x 35 mm
		if err != nil {
			return nil, fmt.Errorf("failed to add barcode to PDF: %v", err)
		}

		// Add product name label
		pdf.SetFont("Arial", "", 8)
		pdf.SetX(xPos)
		pdf.SetY(yPos + 36) // Position below the barcode
		pdf.Cell(nil, product.ProductName)

		// Add price label
		pdf.SetFont("Arial", "", 6)
		pdf.SetX(xPos)
		pdf.SetY(yPos + 43) // Position below the product name
		pdf.Cell(nil, fmt.Sprintf("%s - Rp%.2f", product.ProductID, product.Price))
	}

	// Save the PDF to a file
	var buf bytes.Buffer
	_, err = pdf.WriteTo(&buf)
	if err != nil {
		return nil, errorHandler.NewInternalServer(
			errorHandler.WithInfo("failed to generate PDF"),
			errorHandler.WithMessage(err.Error()),
		)
	}

	// Return the PDF file as a response
	response := &payload.GenerateBarcodeResponse{
		PDFBytes: buf.Bytes(),
	}
	return response, nil
}