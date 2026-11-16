package handler

import (
	"encoding/json"
	"strings"

	"portal_api/domain/auth/dto"
	"portal_api/domain/auth/service"
	global_dto "portal_api/dto"
	"portal_api/errors"
	"portal_api/helper"
	response_helper "portal_api/helper/response"
	"portal_api/pkg/binder"
	validator "portal_api/validation"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthServiceInterface
}

func NewAuthHandler(svc service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{service: svc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	req, err := binder.BindJSON[dto.LoginRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.Login(c.Request.Context(), &req, c.ClientIP())
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Login berhasil",
		Data:    data,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token := extractBearerToken(c)
	if token == "" {
		c.Error(&errors.UnauthenticatedError{Message: "Token tidak ditemukan"})
		return
	}

	if err := h.service.Logout(token); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Logout berhasil",
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	req, err := binder.BindJSON[dto.RefreshRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Token diperbarui",
		Data:    data,
	})
}

// GetMe: profil akun terbaru dari BE User, diteruskan apa adanya.
func (h *AuthHandler) GetMe(c *gin.Context) {
	resp, err := h.service.GetMe(c.Request.Context(), helper.GetAuthUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WriteRaw(c, resp.StatusCode, resp.Body)
}

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

	body, err := json.Marshal(req)
	if err != nil {
		c.Error(err)
		return
	}

	resp, err := h.service.ChangeMyPassword(c.Request.Context(), helper.GetUserID(c), helper.GetAuthUserID(c), body)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WriteRaw(c, resp.StatusCode, resp.Body)
}

func extractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}
