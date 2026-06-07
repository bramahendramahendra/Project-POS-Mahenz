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

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	req, err := binder.BindJSON[dto.ProductListRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	filter := &dto.ProductFilter{
		Search:     req.Search,
		CategoryID: req.CategoryID,
		IsActive:   req.IsActive,
		LowStock:   req.LowStock,
		Page:       req.Page,
		Limit:      req.Limit,
	}

	products, total, svcErr := h.service.GetAll(filter)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar produk",
		Data: gin.H{
			"items": products,
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

// POST /products/search
func (h *ProductHandler) Search(c *gin.Context) {
	req, err := binder.BindJSON[dto.SearchProductRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	results, svcErr := h.service.Search(req.Q, limit)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Hasil pencarian produk",
		Data:    results,
	})
}

// POST /products/by-barcode/:barcode
func (h *ProductHandler) GetByBarcode(c *gin.Context) {
	req, err := binder.BindURI[dto.GetProductByBarcodeRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	product, svcErr := h.service.GetByBarcode(req.Barcode)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Detail produk",
		Data:    product,
	})
}

// POST /products/detail/:id
func (h *ProductHandler) GetByID(c *gin.Context) {
	req, err := binder.BindURI[dto.GetProductByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	product, svcErr := h.service.GetByID(req.ID)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Detail produk",
		Data:    product,
	})
}

// POST /products/create
func (h *ProductHandler) Create(c *gin.Context) {
	req, err := binder.BindJSON[dto.ProductRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	product, svcErr := h.service.Create(&req)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Produk berhasil dibuat",
		Data:    product,
	})
}

// POST /products/update/:id
func (h *ProductHandler) Update(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.UpdateProductUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, bindErr := binder.BindJSON[dto.ProductRequest](c)
	if bindErr != nil {
		c.Error(&errors.BadRequestError{Message: bindErr.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	if svcErr := h.service.Update(uriReq.ID, &req); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Produk berhasil diperbarui",
	})
}

// POST /products/delete/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	req, err := binder.BindURI[dto.DeleteProductRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	if svcErr := h.service.Delete(req.ID); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Produk berhasil dihapus",
	})
}

// POST /products/toggle-status/:id
func (h *ProductHandler) ToggleStatus(c *gin.Context) {
	req, err := binder.BindURI[dto.ToggleStatusProductRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	if svcErr := h.service.ToggleStatus(req.ID); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Status produk berhasil diubah",
	})
}
