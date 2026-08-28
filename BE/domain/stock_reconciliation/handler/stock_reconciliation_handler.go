package handler

import (
	dto "pos_api/domain/stock_reconciliation/dto"
	service "pos_api/domain/stock_reconciliation/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type StockReconciliationHandler struct {
	service service.StockReconciliationServiceInterface
}

func NewStockReconciliationHandler(service service.StockReconciliationServiceInterface) *StockReconciliationHandler {
	return &StockReconciliationHandler{service: service}
}

func (h *StockReconciliationHandler) GetList(c *gin.Context) {
	req, err := binder.BindJSON[dto.ListRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetList(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Daftar rekonsiliasi stok",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *StockReconciliationHandler) GetSummary(c *gin.Context) {
	data, err := h.service.GetSummary()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Ringkasan rekonsiliasi stok",
		Data:    data,
	})
}

func (h *StockReconciliationHandler) GetDetail(c *gin.Context) {
	req, err := binder.BindURI[dto.DetailRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetDetail(req.ProductID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Detail rekonsiliasi stok produk",
		Data:    data,
	})
}

func (h *StockReconciliationHandler) Adjust(c *gin.Context) {
	req, err := binder.BindJSON[dto.AdjustRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	req.UserID = helper.GetUserID(c)

	data, err := h.service.Adjust(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Stok berhasil dikoreksi",
		Data:    data,
	})
}

func (h *StockReconciliationHandler) MarkReviewed(c *gin.Context) {
	req, err := binder.BindJSON[dto.MarkReviewedRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	req.UserID = helper.GetUserID(c)

	data, err := h.service.MarkReviewed(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Produk ditandai sudah ditinjau",
		Data:    data,
	})
}
