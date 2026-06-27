package handler

import (
	dto "pos_api/domain/supplier_return/dto"
	service "pos_api/domain/supplier_return/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	"pos_api/validation"

	"github.com/gin-gonic/gin"
)

type SupplierReturnHandler struct {
	service service.SupplierReturnServiceInterface
}

func NewSupplierReturnHandler(service service.SupplierReturnServiceInterface) *SupplierReturnHandler {
	return &SupplierReturnHandler{service: service}
}

func (h *SupplierReturnHandler) GetAll(c *gin.Context) {
	req, err := binder.BindJSON[dto.SupplierReturnListRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, total, err := h.service.GetAll(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Daftar retur supplier",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *SupplierReturnHandler) GetByID(c *gin.Context) {
	req, err := binder.BindURI[dto.GetSupplierReturnByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetByID(req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Detail retur supplier",
		Data:    data,
	})
}

func (h *SupplierReturnHandler) Create(c *gin.Context) {
	req, err := binder.BindJSON[dto.CreateSupplierReturnRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	req.UserID = helper.GetUserID(c)

	data, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusCreated,
		Status:  true,
		Message: "Retur supplier berhasil dibuat",
		Data:    data,
	})
}

func (h *SupplierReturnHandler) UpdateStatus(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.GetSupplierReturnByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.UpdateStatusRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req.ID = uriReq.ID

	if err := validation.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	req.UserID = helper.GetUserID(c)

	if err := h.service.UpdateStatus(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Status retur supplier berhasil diperbarui",
	})
}

func (h *SupplierReturnHandler) Delete(c *gin.Context) {
	req, err := binder.BindURI[dto.GetSupplierReturnByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validation.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Delete(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Retur supplier berhasil dihapus",
	})
}
