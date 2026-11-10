package handler

import (
	"portal_api/domain/sso/dto"
	"portal_api/domain/sso/service"
	global_dto "portal_api/dto"
	"portal_api/errors"
	"portal_api/helper"
	response_helper "portal_api/helper/response"
	"portal_api/pkg/binder"
	validator "portal_api/validation"

	"github.com/gin-gonic/gin"
)

type SSOHandler struct {
	service service.SSOServiceInterface
}

func NewSSOHandler(svc service.SSOServiceInterface) *SSOHandler {
	return &SSOHandler{service: svc}
}

// Launch: user yang sudah login meminta kode SSO untuk membuka sebuah aplikasi.
func (h *SSOHandler) Launch(c *gin.Context) {
	req, err := binder.BindJSON[dto.LaunchRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.Launch(helper.GetUserID(c), &req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Kode SSO dibuat",
		Data:    data,
	})
}

// Exchange: dipanggil BE aplikasi (server ke server) untuk menukar kode SSO menjadi identitas user.
// Autentikasi memakai kunci aplikasi pada header X-App-Key.
func (h *SSOHandler) Exchange(c *gin.Context) {
	req, err := binder.BindJSON[dto.ExchangeRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.Exchange(c.GetHeader("X-App-Key"), &req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Kode SSO valid",
		Data:    data,
	})
}
