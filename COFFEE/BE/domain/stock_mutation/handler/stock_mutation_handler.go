package handler

import (
	dto "coffee_api/domain/stock_mutation/dto"
	service "coffee_api/domain/stock_mutation/service"
	global_dto "coffee_api/dto"
	"coffee_api/errors"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"
	binder "coffee_api/pkg/binder"
	validator "coffee_api/validation"

	"github.com/gin-gonic/gin"
)

type StockMutationHandler struct {
	service service.StockMutationServiceInterface
}

func NewStockMutationHandler(service service.StockMutationServiceInterface) *StockMutationHandler {
	return &StockMutationHandler{service: service}
}

func (h *StockMutationHandler) GetAll(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetAll(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Daftar mutasi stok",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *StockMutationHandler) GetByProduct(c *gin.Context) {
	req, err := binder.BindURI[dto.GetByProductRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetByProduct(req.ProductID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Riwayat mutasi stok produk",
		Data:    data,
	})
}
