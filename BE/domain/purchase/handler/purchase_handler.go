package handler

import (
	dto "pos_api/domain/purchase/dto"
	service "pos_api/domain/purchase/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	"pos_api/validation"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	service service.PurchaseService
}

func NewPurchaseHandler(service service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: service}
}

func (h *PurchaseHandler) GenerateCode(c *gin.Context) {
	data, err := h.service.GenerateCode()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Generate kode PO",
		Data:    data,
	})
}

func (h *PurchaseHandler) GetAll(c *gin.Context) {
	req, err := binder.BindJSON[dto.PurchaseListRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
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
		Message:    "Daftar purchase order",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, int64(total)),
	})
}

func (h *PurchaseHandler) GetByID(c *gin.Context) {
	req, err := binder.BindURI[dto.GetPurchaseByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
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
		Message: "Detail purchase order",
		Data:    data,
	})
}

func (h *PurchaseHandler) GetItems(c *gin.Context) {
	req, err := binder.BindURI[dto.GetPurchaseByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetItems(req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Item purchase order",
		Data:    data,
	})
}

func (h *PurchaseHandler) Create(c *gin.Context) {
	req, err := binder.BindJSON[dto.PurchaseRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}
	req.UserID = helper.GetUserID(c)

	data, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusCreated,
		Status:  true,
		Message: "Purchase order berhasil dibuat",
		Data:    data,
	})
}

func (h *PurchaseHandler) Update(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.GetPurchaseByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.PurchaseRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.ID

	if err := validation.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.Update(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Purchase order berhasil diperbarui",
		Data:    data,
	})
}

func (h *PurchaseHandler) Delete(c *gin.Context) {
	req, err := binder.BindURI[dto.GetPurchaseByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Delete(req.ID); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Purchase order berhasil dihapus",
	})
}

func (h *PurchaseHandler) GetPayments(c *gin.Context) {
	req, err := binder.BindURI[dto.GetPurchaseByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetPayments(req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Riwayat pembayaran purchase order",
		Data:    data,
	})
}

func (h *PurchaseHandler) Pay(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.GetPurchaseByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.PayPurchaseRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.ID
	req.UserID = helper.GetUserID(c)

	if err := validation.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Pay(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Pembayaran purchase order berhasil dicatat",
	})
}
