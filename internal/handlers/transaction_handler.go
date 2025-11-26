package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"nabung-emas-api/internal/middleware"
	"nabung-emas-api/internal/models"
	"nabung-emas-api/internal/services"
	"nabung-emas-api/internal/utils"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) GetAll(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	// Parse query parameters
	pocketID := c.QueryParam("pocket_id")
	brand := c.QueryParam("brand")
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	sortBy := c.QueryParam("sort_by")
	sortOrder := c.QueryParam("sort_order")

	var pocketIDPtr, brandPtr, startDatePtr, endDatePtr *string
	if pocketID != "" {
		pocketIDPtr = &pocketID
	}
	if brand != "" {
		brandPtr = &brand
	}
	if startDate != "" {
		startDatePtr = &startDate
	}
	if endDate != "" {
		endDatePtr = &endDate
	}

	transactions, total, err := h.service.GetAll(userID, pocketIDPtr, brandPtr, startDatePtr, endDatePtr, page, limit, sortBy, sortOrder)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch transactions")
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	return utils.PaginatedResponse(c, transactions, page, limit, total)
}

func (h *TransactionHandler) GetByID(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	id := c.Param("id")

	transaction, err := h.service.GetByID(id, userID)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
	}

	return utils.SuccessResponse(c, http.StatusOK, "Success", transaction)
}

func (h *TransactionHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	var req models.CreateTransactionRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	transaction, err := h.service.Create(userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusCreated, "Transaction created successfully", transaction)
}

func (h *TransactionHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	id := c.Param("id")

	var req models.UpdateTransactionRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		return utils.HandleError(c, err)
	}

	transaction, err := h.service.Update(id, userID, &req)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Transaction updated successfully", transaction)
}

func (h *TransactionHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	id := c.Param("id")

	if err := h.service.Delete(id, userID); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, http.StatusOK, "Transaction deleted successfully", nil)
}

func (h *TransactionHandler) UploadReceipt(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	id := c.Param("id")

	// TODO: Implement file upload logic
	// For now, return not implemented
	_ = id
	return utils.ErrorResponse(c, http.StatusNotImplemented, "Receipt upload not yet implemented")
}
