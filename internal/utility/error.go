package utility

import "rizkysr90-pos/pkg/errorHandler"

func ConstructErrorRequired(fieldName string) *errorHandler.HttpError {
	return &errorHandler.HttpError{
		Code:    400,
		Info:    fieldName + " is required",
		Message: "",
	}
}
func ConstructErrorMaxLen(fieldName string) *errorHandler.HttpError {
	return &errorHandler.HttpError{
		Code:    400,
		Info:    fieldName + " too long",
		Message: "",
	}
}
