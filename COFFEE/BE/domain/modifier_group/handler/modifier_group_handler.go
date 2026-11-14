package handler

import (
	dto "coffee_api/domain/modifier_group/dto"
	service "coffee_api/domain/modifier_group/service"
	global_dto "coffee_api/dto"
	"coffee_api/errors"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"
	binder "coffee_api/pkg/binder"
	validator "coffee_api/validation"

	"github.com/gin-gonic/gin"
)

type ModifierGroupHandler struct {
	service service.ModifierGroupServiceInterface
}

func NewModifierGroupHandler(service service.ModifierGroupServiceInterface) *ModifierGroupHandler {
	return &ModifierGroupHandler{service: service}
}

func ok(c *gin.Context, status int, message string, data any) {
	response_helper.WrapResponse(c, status, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func bindJSON[T any](c *gin.Context) (T, bool) {
	req, err := binder.BindJSON[T](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return req, false
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return req, false
	}
	return req, true
}

func bindID(c *gin.Context) (int, bool) {
	req, err := binder.BindURI[dto.IDUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return 0, false
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return 0, false
	}
	return req.ID, true
}

func (h *ModifierGroupHandler) GetAll(c *gin.Context) {
	req, valid := bindJSON[dto.GetAllRequest](c)
	if !valid {
		return
	}
	data, err := h.service.GetAll(&req)
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Daftar kelompok tambahan", data)
}

func (h *ModifierGroupHandler) GetOptions(c *gin.Context) {
	data, err := h.service.GetOptions()
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Pilihan kelompok tambahan", data)
}

func (h *ModifierGroupHandler) GetByID(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	data, err := h.service.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Detail kelompok tambahan", data)
}

func (h *ModifierGroupHandler) Create(c *gin.Context) {
	req, valid := bindJSON[dto.SaveRequest](c)
	if !valid {
		return
	}
	data, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 201, "Kelompok tambahan berhasil ditambahkan", data)
}

func (h *ModifierGroupHandler) Update(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	req, valid := bindJSON[dto.SaveRequest](c)
	if !valid {
		return
	}
	req.ID = id
	if err := h.service.Update(&req); err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Kelompok tambahan berhasil diperbarui", nil)
}

func (h *ModifierGroupHandler) ToggleStatus(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	if err := h.service.ToggleStatus(id); err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Status kelompok tambahan berhasil diubah", nil)
}

func (h *ModifierGroupHandler) Delete(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	if err := h.service.Delete(id); err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Kelompok tambahan berhasil dihapus", nil)
}
