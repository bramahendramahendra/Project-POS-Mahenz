package handler

import (
	dto "coffee_api/domain/menu_item/dto"
	service "coffee_api/domain/menu_item/service"
	global_dto "coffee_api/dto"
	"coffee_api/errors"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"
	binder "coffee_api/pkg/binder"
	validator "coffee_api/validation"

	"github.com/gin-gonic/gin"
)

type MenuItemHandler struct {
	service service.MenuItemServiceInterface
}

func NewMenuItemHandler(service service.MenuItemServiceInterface) *MenuItemHandler {
	return &MenuItemHandler{service: service}
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

// Catalog: data halaman kasir (tahap 4.7b).
func (h *MenuItemHandler) Catalog(c *gin.Context) {
	data, err := h.service.Catalog()
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Katalog kasir", data)
}

func (h *MenuItemHandler) GetAll(c *gin.Context) {
	req, valid := bindJSON[dto.GetAllRequest](c)
	if !valid {
		return
	}
	data, total, err := h.service.GetAll(&req)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Daftar menu",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *MenuItemHandler) GetByID(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	data, err := h.service.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Detail menu", data)
}

func (h *MenuItemHandler) Create(c *gin.Context) {
	req, valid := bindJSON[dto.SaveRequest](c)
	if !valid {
		return
	}
	data, err := h.service.Create(&req, helper.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 201, "Menu berhasil ditambahkan", data)
}

func (h *MenuItemHandler) Update(c *gin.Context) {
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
	ok(c, 200, "Menu berhasil diperbarui", nil)
}

func (h *MenuItemHandler) ToggleStatus(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	if err := h.service.ToggleStatus(id); err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Status menu berhasil diubah", nil)
}

func (h *MenuItemHandler) SetSoldOut(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	req, valid := bindJSON[dto.SoldOutRequest](c)
	if !valid {
		return
	}
	req.ID = id
	if err := h.service.SetSoldOut(&req); err != nil {
		c.Error(err)
		return
	}
	message := "Menu ditandai tersedia"
	if *req.IsSoldOut {
		message = "Menu ditandai habis"
	}
	ok(c, 200, message, nil)
}

func (h *MenuItemHandler) Delete(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	if err := h.service.Delete(id); err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Menu berhasil dihapus", nil)
}
