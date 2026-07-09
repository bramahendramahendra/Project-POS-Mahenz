package handler

import (
	dto "pos_api/domain/finance/dto"
	service "pos_api/domain/finance/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type FinanceHandler struct {
	service service.FinanceServiceInterface
}

func NewFinanceHandler(service service.FinanceServiceInterface) *FinanceHandler {
	return &FinanceHandler{service: service}
}

func (h *FinanceHandler) GetSummary(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetSummaryRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.GetSummary(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Ringkasan keuangan",
		Data:    data,
	})
}

func (h *FinanceHandler) GetCashflow(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetCashflowRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetCashflow(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Daftar arus kas",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}
