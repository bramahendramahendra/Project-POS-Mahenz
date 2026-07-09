package handler

import (
	"pos_api/domain/setting/dto"
	"pos_api/domain/setting/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	service service.SettingServiceInterface
}

func NewSettingHandler(svc service.SettingServiceInterface) *SettingHandler {
	return &SettingHandler{service: svc}
}

// GET /api/settings
func (h *SettingHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll()
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar pengaturan",
		Data:    data,
	})
}

// GET /api/settings/:key
func (h *SettingHandler) GetByKey(c *gin.Context) {
	key := c.Param("key")
	value, err := h.service.GetByKey(key)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Detail pengaturan",
		Data:    dto.SettingKeyValueResponse{Key: key, Value: value},
	})
}

// POST /api/settings
func (h *SettingHandler) Save(c *gin.Context) {
	body, err := binder.BindJSON[map[string]string](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := h.service.Save(body); err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Settings berhasil disimpan",
	})
}

// POST /api/settings/reset
func (h *SettingHandler) Reset(c *gin.Context) {
	if err := h.service.Reset(); err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Settings berhasil direset ke default",
	})
}

// GET /api/settings/store
func (h *SettingHandler) GetStoreProfile(c *gin.Context) {
	data, err := h.service.GetStoreProfile()
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Profil toko",
		Data:    data,
	})
}

// POST /api/settings/store
func (h *SettingHandler) UpdateStoreProfile(c *gin.Context) {
	req, err := binder.BindJSON[dto.StoreProfileRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := h.service.UpdateStoreProfile(&req); err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Profil toko berhasil disimpan",
	})
}

// GET /api/settings/printer
func (h *SettingHandler) GetPrinterSettings(c *gin.Context) {
	data, err := h.service.GetPrinterSettings()
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Pengaturan printer",
		Data:    data,
	})
}

// POST /api/settings/printer
func (h *SettingHandler) UpdatePrinterSettings(c *gin.Context) {
	req, err := binder.BindJSON[dto.PrinterSettingsRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := h.service.UpdatePrinterSettings(&req); err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Pengaturan printer berhasil disimpan",
	})
}
