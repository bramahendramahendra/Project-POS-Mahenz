package handler

import (
	dto_purchase "pos_api/domain/purchase/dto"
	service_purchase "pos_api/domain/purchase/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	"pos_api/validation"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	service service_purchase.PurchaseService
}

func NewPurchaseHandler(service service_purchase.PurchaseService) *PurchaseHandler {
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
	req, err := binder.BindJSON[dto_purchase.PurchaseListRequest](c)
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
	req, err := binder.BindURI[dto_purchase.GetPurchaseByIDRequest](c)
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
	req, err := binder.BindURI[dto_purchase.GetPurchaseByIDRequest](c)
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
	req, err := binder.BindJSON[dto_purchase.PurchaseRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	data, err := h.service.Create(&req, uid)
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

// POST /supplier-purchases/update/:id
func (h *PurchaseHandler) Update(c *gin.Context) {
	uriReq, err := binder.BindURI[dto_purchase.GetPurchaseByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto_purchase.PurchaseRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.ID

	if err := validation.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, svcErr := h.service.Update(req.ID, &req)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Purchase order berhasil diperbarui",
		Data:    data,
	})
}

// POST /supplier-purchases/delete/:id
func (h *PurchaseHandler) Delete(c *gin.Context) {
	req, err := binder.BindURI[dto_purchase.GetPurchaseByIDRequest](c)
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
	req, err := binder.BindURI[dto_purchase.GetPurchaseByIDRequest](c)
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

// POST /supplier-purchases/pay/:id
func (h *PurchaseHandler) Pay(c *gin.Context) {
	uriReq, err := binder.BindURI[dto_purchase.GetPurchaseByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto_purchase.PayPurchaseRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.ID

	if err := validation.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	if svcErr := h.service.Pay(req.ID, &req, uid); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Pembayaran purchase order berhasil dicatat",
	})
}
