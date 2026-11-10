package handler

import (
	dto "portal_api/domain/user/dto"
	service "portal_api/domain/user/service"
	global_dto "portal_api/dto"
	"portal_api/errors"
	"portal_api/helper"
	response_helper "portal_api/helper/response"
	binder "portal_api/pkg/binder"
	validator "portal_api/validation"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, svcErr := h.service.GetAll(&req)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Daftar user",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
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
		Message: "Detail user",
		Data:    data,
	})
}

func (h *UserHandler) Create(c *gin.Context) {
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
		Message: "User berhasil dibuat",
		Data:    data,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
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

	if svcErr := h.service.Update(uriReq.ID, helper.GetUserID(c), &bodyReq); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "User berhasil diperbarui",
	})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.ResetPasswordUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(uriReq); err != nil {
		c.Error(err)
		return
	}

	bodyReq, err := binder.BindJSON[dto.ResetPasswordRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(bodyReq); err != nil {
		c.Error(err)
		return
	}

	if svcErr := h.service.ResetPassword(uriReq.ID, helper.GetUserID(c), &bodyReq); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Password sementara berhasil ditetapkan. User wajib menggantinya saat login berikutnya",
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	req, err := binder.BindURI[dto.DeleteRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if svcErr := h.service.Delete(req.ID, helper.GetUserID(c)); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "User berhasil dihapus",
	})
}

func (h *UserHandler) ToggleStatus(c *gin.Context) {
	req, err := binder.BindURI[dto.ToggleStatusRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if svcErr := h.service.ToggleStatus(req.ID, helper.GetUserID(c)); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Status user berhasil diubah",
	})
}
