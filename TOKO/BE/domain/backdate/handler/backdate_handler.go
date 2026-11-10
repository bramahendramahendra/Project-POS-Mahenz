package handler

import (
	"pos_api/domain/backdate/dto"
	"pos_api/domain/backdate/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type BackdateHandler struct {
	service service.BackdateServiceInterface
}

func NewBackdateHandler(service service.BackdateServiceInterface) *BackdateHandler {
	return &BackdateHandler{service: service}
}

func (h *BackdateHandler) OpenCashDrawer(c *gin.Context) {
	req, err := binder.BindJSON[dto.OpenCashDrawerRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	adminUserID := helper.GetUserID(c)

	data, err := h.service.OpenCashDrawer(adminUserID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Kas historis berhasil dibuka",
		Data:    data,
	})
}

func (h *BackdateHandler) GetCurrentCashDrawer(c *gin.Context) {
	adminUserID := helper.GetUserID(c)

	data, err := h.service.GetCurrentCashDrawer(adminUserID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Status kas historis",
		Data:    data,
	})
}

func (h *BackdateHandler) CloseCashDrawer(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.CloseCashDrawerRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	type closeBody struct {
		ClosingBalance float64 `json:"closing_balance" validate:"min=0"`
		Notes          string  `json:"notes"`
	}
	body, err := binder.BindJSON[closeBody](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	uriReq.ClosingBalance = body.ClosingBalance
	uriReq.Notes = body.Notes

	adminUserID := helper.GetUserID(c)

	data, err := h.service.CloseCashDrawer(adminUserID, &uriReq)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Kas historis berhasil ditutup",
		Data:    data,
	})
}

func (h *BackdateHandler) CreateTransaction(c *gin.Context) {
	req, err := binder.BindJSON[dto.CreateTransactionRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	adminUserID := helper.GetUserID(c)

	data, err := h.service.CreateTransaction(adminUserID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Transaksi historis berhasil dibuat",
		Data:    data,
	})
}
