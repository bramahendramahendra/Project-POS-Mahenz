package handler

import (
	"pos_api/domain/version/dto"
	"pos_api/domain/version/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	"pos_api/pkg/binder"

	"github.com/gin-gonic/gin"
)

type VersionHandler struct {
	service service.VersionServiceInterface
}

func NewVersionHandler(svc service.VersionServiceInterface) *VersionHandler {
	return &VersionHandler{service: svc}
}

// GET /api/version/android
func (h *VersionHandler) CheckAndroid(c *gin.Context) {
	currentVersion := c.Query("current_version")
	if currentVersion == "" {
		c.Error(&errors.BadRequestError{Message: "current_version wajib diisi"})
		return
	}

	result, err := h.service.CheckAndroid(currentVersion)
	if err != nil {
		c.Error(err)
		return
	}

	msg := "Aplikasi sudah versi terbaru"
	if result.HasUpdate {
		msg = "Versi baru tersedia"
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: msg,
		Data:    result,
	})
}

// POST /api/version/android
func (h *VersionHandler) UpdateAndroidVersion(c *gin.Context) {
	req, err := binder.BindJSON[dto.UpdateVersionRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := h.service.UpdateAndroidVersion(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusCreated,
		Status:  true,
		Message: "Versi berhasil diupdate",
	})
}
