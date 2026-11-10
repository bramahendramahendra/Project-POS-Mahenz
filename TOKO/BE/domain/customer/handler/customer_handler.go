package handler

import (
	dto "pos_api/domain/customer/dto"
	service "pos_api/domain/customer/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service service.CustomerServiceInterface
}

func NewCustomerHandler(service service.CustomerServiceInterface) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) GetAll(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
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
		Message:    "Daftar pelanggan",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *CustomerHandler) GetOptions(c *gin.Context) {
	data, err := h.service.GetOptions()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Opsi pelanggan",
		Data:    data,
	})
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	req, err := binder.BindURI[dto.GetByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
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
		Message: "Detail pelanggan",
		Data:    data,
	})
}

func (h *CustomerHandler) Create(c *gin.Context) {
	req, err := binder.BindJSON[dto.CreateRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Pelanggan berhasil ditambahkan",
		Data:    data,
	})
}

func (h *CustomerHandler) Update(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.UpdateUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.UpdateRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.ID

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.Update(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Pelanggan berhasil diperbarui",
		Data:    data,
	})
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	req, err := binder.BindURI[dto.DeleteRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
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
		Message: "Pelanggan berhasil dihapus",
	})
}

func (h *CustomerHandler) ToggleStatus(c *gin.Context) {
	req, err := binder.BindURI[dto.ToggleStatusRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.ToggleStatus(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Status pelanggan berhasil diubah",
	})
}
