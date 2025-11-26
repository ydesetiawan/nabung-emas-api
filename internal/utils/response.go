package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"nabung-emas-api/internal/models"
)

func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, models.APIResponse{
		Success: false,
		Message: message,
	})
}

func ValidationErrorResponse(c echo.Context, err error) error {
	errors := GetValidationErrors(err)
	return c.JSON(http.StatusUnprocessableEntity, models.APIResponse{
		Success: false,
		Message: "Validation failed",
		Errors:  errors,
	})
}

func PaginatedResponse(c echo.Context, data interface{}, page, limit, total int) error {
	totalPages := (total + limit - 1) / limit
	
	return c.JSON(http.StatusOK, models.PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: models.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func HandleError(c echo.Context, err error) error {
	if he, ok := err.(*echo.HTTPError); ok {
		return ErrorResponse(c, he.Code, he.Message.(string))
	}

	if _, ok := err.(validator.ValidationErrors); ok {
		return ValidationErrorResponse(c, err)
	}

	return ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
}
