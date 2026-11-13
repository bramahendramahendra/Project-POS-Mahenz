package handler

import (
	"auth_api/domain/auth/dto"
	"auth_api/domain/auth/service"
	global_dto "auth_api/dto"
	"auth_api/errors"
	"auth_api/helper"
	response_helper "auth_api/helper/response"
	"auth_api/pkg/binder"
	validator "auth_api/validation"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthServiceInterface
}

func NewAuthHandler(svc service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{service: svc}
}

// Verify: dipanggil BE aplikasi (Portal, Toko, Coffee) untuk memeriksa username + password.
func (h *AuthHandler) Verify(c *gin.Context) {
	req, err := binder.BindJSON[dto.VerifyRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.Verify(&req, helper.GetCallerApp(c))
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Login valid",
		Data:    data,
	})
}

// Me: profil akun yang bertindak (dipanggil Portal atas nama user).
func (h *AuthHandler) Me(c *gin.Context) {
	data, err := h.service.GetMe(helper.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Data user",
		Data:    data,
	})
}

// ChangePassword: ganti password milik sendiri (dipanggil Portal atas nama user).
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	req, err := binder.BindJSON[dto.ChangeMyPasswordRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if err := h.service.ChangeMyPassword(helper.GetUserID(c), &req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Password berhasil diubah",
	})
}
