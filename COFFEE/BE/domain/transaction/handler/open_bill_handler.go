package handler

import (
	"coffee_api/domain/transaction/dto"
	global_dto "coffee_api/dto"
	"coffee_api/errors"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"
	binder "coffee_api/pkg/binder"
	validator "coffee_api/validation"

	"github.com/gin-gonic/gin"
)

// Open bill (tahap 4.8). Validasi struktur ditolak sebagai 400 seperti handler transaksi lainnya.

func respond(c *gin.Context, status int, message string, data any) {
	response_helper.WrapResponse(c, status, "json", &global_dto.ResponseParams{
		Code: helper.StatusOk, Status: true, Message: message, Data: data,
	})
}

func bindOpenID(c *gin.Context) (int, bool) {
	req, err := binder.BindURI[dto.OpenBillURIRequest](c)
	if err == nil {
		err = validator.Validate.Struct(req)
	}
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return 0, false
	}
	return req.ID, true
}

func bindOpenBody[T any](c *gin.Context) (T, bool) {
	req, err := binder.BindJSON[T](c)
	if err == nil {
		err = validator.Validate.Struct(req)
	}
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return req, false
	}
	return req, true
}

func (h *TransactionHandler) ListOpen(c *gin.Context) {
	data, err := h.service.ListOpen()
	if err != nil {
		c.Error(err)
		return
	}
	respond(c, 200, "Pesanan terbuka", data)
}

func (h *TransactionHandler) AddOpenItems(c *gin.Context) {
	id, ok := bindOpenID(c)
	if !ok {
		return
	}
	req, ok := bindOpenBody[dto.OpenAddItemsRequest](c)
	if !ok {
		return
	}
	data, err := h.service.AddOpenItems(id, &req, helper.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	respond(c, 200, "Item ditambahkan ke pesanan", data)
}

func (h *TransactionHandler) RemoveOpenItem(c *gin.Context) {
	id, ok := bindOpenID(c)
	if !ok {
		return
	}
	req, ok := bindOpenBody[dto.OpenRemoveItemRequest](c)
	if !ok {
		return
	}
	data, err := h.service.RemoveOpenItem(id, &req, helper.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	respond(c, 200, "Item dihapus dari pesanan", data)
}

func (h *TransactionHandler) MoveOpenTable(c *gin.Context) {
	id, ok := bindOpenID(c)
	if !ok {
		return
	}
	req, ok := bindOpenBody[dto.OpenMoveTableRequest](c)
	if !ok {
		return
	}
	data, err := h.service.MoveOpenTable(id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	respond(c, 200, "Pesanan dipindah meja", data)
}

func (h *TransactionHandler) MergeOpen(c *gin.Context) {
	id, ok := bindOpenID(c)
	if !ok {
		return
	}
	req, ok := bindOpenBody[dto.OpenMergeRequest](c)
	if !ok {
		return
	}
	data, err := h.service.MergeOpen(id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	respond(c, 200, "Pesanan digabung", data)
}

func (h *TransactionHandler) CancelOpen(c *gin.Context) {
	id, ok := bindOpenID(c)
	if !ok {
		return
	}
	req, ok := bindOpenBody[dto.OpenCancelRequest](c)
	if !ok {
		return
	}
	if err := h.service.CancelOpen(id, &req, helper.GetUserID(c)); err != nil {
		c.Error(err)
		return
	}
	respond(c, 200, "Pesanan dibatalkan", nil)
}

func (h *TransactionHandler) PayOpen(c *gin.Context) {
	id, ok := bindOpenID(c)
	if !ok {
		return
	}
	req, ok := bindOpenBody[dto.OpenPayRequest](c)
	if !ok {
		return
	}
	data, err := h.service.PayOpen(id, &req, helper.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	respond(c, 200, "Pesanan dibayar", data)
}
