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

type ProductPackageHandler struct {
	service service.ProductServiceInterface
}

func NewProductPackageHandler(service service.ProductServiceInterface) *ProductPackageHandler {
	return &ProductPackageHandler{service: service}
}

// POST /products/:id/packages/list
func (h *ProductPackageHandler) GetPackagesByProduct(c *gin.Context) {
	req, err := binder.BindURI[dto.PackageByProductRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetPackagesByProduct(req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar paket produk",
		Data:    data,
	})
}

// POST /products/:id/packages/create
func (h *ProductPackageHandler) CreatePackage(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.PackageByProductRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.CreatePackageRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ProductID = uriReq.ID

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.CreatePackage(uriReq.ID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Paket produk berhasil ditambahkan",
		Data:    data,
	})
}

// POST /products/:id/packages/update/:package_id
func (h *ProductPackageHandler) UpdatePackage(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.PackageUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.UpdatePackageRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.PackageID
	req.ProductID = uriReq.ID

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.UpdatePackage(uriReq.PackageID, uriReq.ID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Paket produk berhasil diperbarui",
		Data:    data,
	})
}

// POST /products/:id/packages/delete/:package_id
func (h *ProductPackageHandler) DeletePackage(c *gin.Context) {
	req, err := binder.BindURI[dto.DeletePackageRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.DeletePackage(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Paket produk berhasil dihapus",
	})
}
