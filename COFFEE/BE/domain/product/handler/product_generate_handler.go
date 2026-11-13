package handler

import (
	dto "coffee_api/domain/product/dto"
	service "coffee_api/domain/product/service"
	global_dto "coffee_api/dto"
	"coffee_api/errors"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"
	binder "coffee_api/pkg/binder"
	validator "coffee_api/validation"

	"github.com/gin-gonic/gin"
)

type ProductGenerateHandler struct {
	service service.ProductServiceInterface
}

func NewProductGenerateHandler(service service.ProductServiceInterface) *ProductGenerateHandler {
	return &ProductGenerateHandler{service: service}
}

// POST /products/generate-barcode
func (h *ProductGenerateHandler) GenerateBarcode(c *gin.Context) {
	data, err := h.service.GenerateBarcode()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Barcode berhasil digenerate",
		Data:    data,
	})
}

// POST /products/generate-sku
func (h *ProductGenerateHandler) GenerateSku(c *gin.Context) {
	req, err := binder.BindJSON[dto.GenerateSkuRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GenerateSku(req.CategoryID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "SKU berhasil digenerate",
		Data:    data,
	})
}
