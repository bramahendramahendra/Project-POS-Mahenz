package handler

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

type ProductPackageHandler struct {
	service service.ProductServiceInterface
}

func NewProductPackageHandler(service service.ProductServiceInterface) *ProductPackageHandler {
	return &ProductPackageHandler{service: service}
}

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

func (h *ProductPackageHandler) SavePackages(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.PackageByProductRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.SavePackageRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ProductID = uriReq.ID

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.SavePackages(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Paket produk berhasil disimpan",
	})
}

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
