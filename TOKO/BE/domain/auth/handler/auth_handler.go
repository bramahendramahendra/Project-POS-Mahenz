package handler

import (
	"strings"

	"pos_api/domain/auth/dto"
	"pos_api/domain/auth/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	"pos_api/pkg/binder"

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

	data, err := h.service.Login(&req, c.ClientIP())
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

	data, err := h.service.RefreshToken(req.RefreshToken)
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

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := helper.GetUserID(c)

	data, err := h.service.GetMe(userID)
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

func (h *AuthHandler) VerifyToken(c *gin.Context) {
	req, err := binder.BindJSON[dto.VerifyTokenRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.VerifyToken(req.Token)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Hasil verifikasi token",
		Data:    data,
	})
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
