package payload

type GenerateBarcodeRequest struct {
    ProductID string `json:"product_id"`
}
type GenerateBarcodeResponse struct {
    PDFBytes []byte `json:"pdf_bytes"`
}
