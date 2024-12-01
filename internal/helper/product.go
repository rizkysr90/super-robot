package helper

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	productPrefix = "PRD"
	idLength      = 13 // Total length including prefix
)

func GenerateProductID() (string, error) {
	// Generate timestamp component (6 digits: YYMMDD)
	timestamp := time.Now().Format("060102")

	// Generate random component (3 digits)
	randomNum, err := generateRandomNumber(3)
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %w", err)
	}

	// Combine parts
	baseID := fmt.Sprintf("%s%s%03d", productPrefix, timestamp, randomNum)

	// Calculate and append check digit
	checkDigit := calculateCheckDigit(baseID)
	productID := baseID + checkDigit

	return productID, nil
}

func generateRandomNumber(digits int) (int, error) {
	max := big.NewInt(int64(pow(10, digits)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func calculateCheckDigit(input string) string {
	// Remove prefix for calculation
	numericPart := strings.TrimPrefix(input, productPrefix)

	sum := 0
	for i, ch := range numericPart {
		digit := int(ch - '0')
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}
	checkDigit := (10 - (sum % 10)) % 10
	return strconv.Itoa(checkDigit)
}

func pow(base, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
