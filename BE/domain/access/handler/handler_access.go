package handler

import (
	dto "pos_api/domain/access/dto"
	service "pos_api/domain/access/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type AccessHandler struct {
	service service.AccessServiceInterface
}

func NewAccessHandler(service service.AccessServiceInterface) *AccessHandler {
	return &AccessHandler{service: service}
}

func (h *AccessHandler) GetByRoleID(c *gin.Context) {
	req, err := binder.BindURI[dto.GetByRoleIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetByRoleID(req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Akses menu role",
		Data:    data,
	})
}

func (h *AccessHandler) SetRoleAccess(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.SetRoleAccessUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.SetRoleAccessRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.RoleID = uriReq.ID

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.SetRoleAccess(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Akses menu role berhasil disimpan",
	})
}
