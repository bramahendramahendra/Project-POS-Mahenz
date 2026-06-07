package handler_product

import (
	dto_product "pos_api/domain/product/dto"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

// POST /products/generate-barcode
func (h *ProductHandler) GenerateBarcode(c *gin.Context) {
	result, err := h.service.GenerateBarcode()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Barcode berhasil digenerate",
		Data:    result,
	})
}

// POST /products/generate-sku
func (h *ProductHandler) GenerateSku(c *gin.Context) {
	req, err := binder.BindJSON[dto_product.GenerateSkuRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	result, svcErr := h.service.GenerateSku(req.CategoryID)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "SKU berhasil digenerate",
		Data:    result,
	})
}
