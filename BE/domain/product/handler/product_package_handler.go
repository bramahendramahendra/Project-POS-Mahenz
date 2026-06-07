package handler_product

import (
	dto_product "pos_api/domain/product/dto"
	service_product "pos_api/domain/product/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type ProductPackageHandler struct {
	service service_product.ProductPackageService
}

func NewProductPackageHandler(service service_product.ProductPackageService) *ProductPackageHandler {
	return &ProductPackageHandler{service: service}
}

type productIDUriRequest struct {
	ID int `uri:"id" validate:"required,min=1"`
}

type packageIDUriRequest struct {
	ID        int `uri:"id" validate:"required,min=1"`
	PackageID int `uri:"package_id" validate:"required,min=1"`
}

// POST /products/:id/packages/list
func (h *ProductPackageHandler) GetByProduct(c *gin.Context) {
	req, err := binder.BindURI[productIDUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	packages, svcErr := h.service.GetByProduct(req.ID)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar paket produk",
		Data:    packages,
	})
}

// POST /products/:id/packages/save
func (h *ProductPackageHandler) Save(c *gin.Context) {
	uriReq, err := binder.BindURI[productIDUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, bindErr := binder.BindJSON[dto_product.SaveProductPackagesRequest](c)
	if bindErr != nil {
		c.Error(&errors.BadRequestError{Message: bindErr.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	if svcErr := h.service.Save(uriReq.ID, req.Packages); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Paket produk berhasil disimpan",
	})
}

// POST /products/:id/packages/delete/:package_id
func (h *ProductPackageHandler) Delete(c *gin.Context) {
	req, err := binder.BindURI[packageIDUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if valErr := validator.Validate.Struct(req); valErr != nil {
		c.Error(&errors.BadRequestError{Message: valErr.Error()})
		return
	}

	if svcErr := h.service.DeleteOne(req.PackageID, req.ID); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Paket produk berhasil dihapus",
	})
}
