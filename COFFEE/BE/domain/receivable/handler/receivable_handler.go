package handler

import (
	dto "coffee_api/domain/receivable/dto"
	service "coffee_api/domain/receivable/service"
	global_dto "coffee_api/dto"
	"coffee_api/errors"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"
	binder "coffee_api/pkg/binder"
	validator "coffee_api/validation"

	"github.com/gin-gonic/gin"
)

type ReceivableHandler struct {
	service service.ReceivableServiceInterface
}

func NewReceivableHandler(service service.ReceivableServiceInterface) *ReceivableHandler {
	return &ReceivableHandler{service: service}
}

func (h *ReceivableHandler) GetAll(c *gin.Context) {
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
		Message:    "Daftar piutang",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *ReceivableHandler) GetStats(c *gin.Context) {
	data, err := h.service.GetStats()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Ringkasan piutang",
		Data:    data,
	})
}

func (h *ReceivableHandler) GetByID(c *gin.Context) {
	req, err := binder.BindURI[dto.GetByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetByID(req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Detail piutang",
		Data:    data,
	})
}

func (h *ReceivableHandler) Pay(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.PayUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.PayRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req.ID = uriReq.ID

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	req.UserID = helper.GetUserID(c)

	data, err := h.service.Pay(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Pembayaran piutang berhasil",
		Data:    data,
	})
}
