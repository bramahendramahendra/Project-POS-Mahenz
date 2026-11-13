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

type ProductPriceHandler struct {
	service service.ProductServiceInterface
}

func NewProductPriceHandler(service service.ProductServiceInterface) *ProductPriceHandler {
	return &ProductPriceHandler{service: service}
}

// POST /products/:id/prices/list
func (h *ProductPriceHandler) GetPricesByProduct(c *gin.Context) {
	req, err := binder.BindURI[dto.PriceByProductRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetPricesByProduct(req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar harga tier produk",
		Data:    data,
	})
}
