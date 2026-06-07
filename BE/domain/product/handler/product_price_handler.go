package handler_product

import (
	dto "pos_api/domain/product/dto"
	service "pos_api/domain/product/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type ProductPriceHandler struct {
	service service.ProductPriceService
}

func NewProductPriceHandler(service service.ProductPriceService) *ProductPriceHandler {
	return &ProductPriceHandler{service: service}
}

// POST /products/:id/prices/list
func (h *ProductPriceHandler) GetByProduct(c *gin.Context) {
	req, err := binder.BindURI[productIDUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	prices, svcErr := h.service.GetByProduct(req.ID)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar harga tier produk",
		Data:    prices,
	})
}

// POST /products/:id/prices/save
func (h *ProductPriceHandler) Save(c *gin.Context) {
	uriReq, err := binder.BindURI[productIDUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, bindErr := binder.BindJSON[dto.SaveProductPricesRequest](c)
	if bindErr != nil {
		c.Error(&errors.BadRequestError{Message: bindErr.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	if svcErr := h.service.Save(uriReq.ID, req.Prices); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Harga tier produk berhasil disimpan",
	})
}
