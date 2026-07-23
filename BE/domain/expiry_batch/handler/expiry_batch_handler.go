package handler

import (
	dto "pos_api/domain/expiry_batch/dto"
	service "pos_api/domain/expiry_batch/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type ExpiryBatchHandler struct {
	service service.ExpiryBatchServiceInterface
}

func NewExpiryBatchHandler(service service.ExpiryBatchServiceInterface) *ExpiryBatchHandler {
	return &ExpiryBatchHandler{service: service}
}

func (h *ExpiryBatchHandler) GetWarnings(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetWarningsRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, productSeverity, err := h.service.GetWarnings(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar warning expired",
		Data: gin.H{
			"warnings":         data,
			"product_severity": productSeverity,
		},
	})
}

func (h *ExpiryBatchHandler) GetByProduct(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.GetByProductUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(uriReq); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetByProduct(uriReq.ID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Histori batch expired produk",
		Data:    data,
	})
}

func (h *ExpiryBatchHandler) Confirm(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.ConfirmUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.ConfirmRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.ID
	req.UserID = helper.GetUserID(c)

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Confirm(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Batch dikonfirmasi aman, warning dihapus",
	})
}

func (h *ExpiryBatchHandler) WriteOff(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.WriteOffUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.WriteOffRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.ID
	req.UserID = helper.GetUserID(c)

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.WriteOff(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Stok expired berhasil dimusnahkan",
	})
}
