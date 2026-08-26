package handler

import (
	dto "pos_api/domain/product/dto"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"

	"github.com/gin-gonic/gin"
)

// POST /products/:id/purchase-history
func (h *ProductHandler) GetPurchaseHistory(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.ProductHistoryRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	bodyReq, err := binder.BindJSON[dto.ProductHistoryRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	page := bodyReq.Page
	limit := bodyReq.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	data, err := h.service.GetPurchaseHistory(uriReq.ID, page, limit)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Riwayat pembelian produk",
		Data:    data,
	})
}

// POST /products/:id/sale-history
func (h *ProductHandler) GetSaleHistory(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.ProductHistoryRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	bodyReq, err := binder.BindJSON[dto.ProductHistoryRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	page := bodyReq.Page
	limit := bodyReq.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	data, err := h.service.GetSaleHistory(uriReq.ID, page, limit, bodyReq.Status)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Riwayat penjualan produk",
		Data:    data,
	})
}
