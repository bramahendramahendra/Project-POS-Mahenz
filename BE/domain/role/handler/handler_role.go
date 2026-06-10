package handler

import (
	dto "pos_api/domain/role/dto"
	service "pos_api/domain/role/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service service.RoleServiceInterface
}

func NewRoleHandler(service service.RoleServiceInterface) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) GetAll(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, svcErr := h.service.GetAll(&req)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar role",
		Data:    data,
	})
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	req, err := binder.BindURI[dto.GetByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, svcErr := h.service.GetByID(req.ID)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Detail role",
		Data:    data,
	})
}

func (h *RoleHandler) Create(c *gin.Context) {
	req, err := binder.BindJSON[dto.CreateRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, svcErr := h.service.Create(&req)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Role berhasil ditambahkan",
		Data:    data,
	})
}

func (h *RoleHandler) Update(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.UpdateUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(uriReq); err != nil {
		c.Error(err)
		return
	}

	bodyReq, err := binder.BindJSON[dto.UpdateRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(bodyReq); err != nil {
		c.Error(err)
		return
	}

	if svcErr := h.service.Update(uriReq.ID, &bodyReq); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Role berhasil diperbarui",
	})
}

func (h *RoleHandler) Delete(c *gin.Context) {
	req, err := binder.BindURI[dto.DeleteRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if svcErr := h.service.Delete(req.ID); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Role berhasil dihapus",
	})
}

func (h *RoleHandler) ToggleStatus(c *gin.Context) {
	req, err := binder.BindURI[dto.ToggleStatusRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if svcErr := h.service.ToggleStatus(req.ID); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Status role berhasil diubah",
	})
}
