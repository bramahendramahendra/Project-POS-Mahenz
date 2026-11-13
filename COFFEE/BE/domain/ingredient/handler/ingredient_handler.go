package handler

import (
	dto "coffee_api/domain/ingredient/dto"
	service "coffee_api/domain/ingredient/service"
	global_dto "coffee_api/dto"
	"coffee_api/errors"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"
	binder "coffee_api/pkg/binder"
	validator "coffee_api/validation"

	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	service service.IngredientServiceInterface
}

func NewIngredientHandler(service service.IngredientServiceInterface) *IngredientHandler {
	return &IngredientHandler{service: service}
}

func ok(c *gin.Context, status int, message string, data any) {
	response_helper.WrapResponse(c, status, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// bindJSON membaca dan memvalidasi body. false = galat sudah dicatat ke context.
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

func (h *IngredientHandler) GetAll(c *gin.Context) {
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
		Message:    "Daftar bahan",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *IngredientHandler) GetOptions(c *gin.Context) {
	data, err := h.service.GetOptions()
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Pilihan bahan", data)
}

func (h *IngredientHandler) GetByID(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	data, err := h.service.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Detail bahan", data)
}

func (h *IngredientHandler) Create(c *gin.Context) {
	req, valid := bindJSON[dto.CreateRequest](c)
	if !valid {
		return
	}
	data, err := h.service.Create(&req, helper.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 201, "Bahan berhasil ditambahkan", data)
}

func (h *IngredientHandler) Update(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	req, valid := bindJSON[dto.UpdateRequest](c)
	if !valid {
		return
	}
	req.ID = id
	if err := h.service.Update(&req); err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Bahan berhasil diperbarui", nil)
}

func (h *IngredientHandler) ToggleStatus(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	if err := h.service.ToggleStatus(id); err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Status bahan berhasil diubah", nil)
}

func (h *IngredientHandler) Delete(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	if err := h.service.Delete(id); err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Bahan berhasil dihapus", nil)
}

func (h *IngredientHandler) RecordUsage(c *gin.Context) {
	req, valid := bindJSON[dto.RecordUsageRequest](c)
	if !valid {
		return
	}
	data, err := h.service.RecordUsage(&req, helper.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 201, recordMessage("Pemakaian bahan dicatat", data), data)
}

func (h *IngredientHandler) RecordWaste(c *gin.Context) {
	req, valid := bindJSON[dto.RecordWasteRequest](c)
	if !valid {
		return
	}
	data, err := h.service.RecordWaste(&req, helper.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 201, recordMessage("Barang terbuang dicatat", data), data)
}

// recordMessage menambahkan peringatan bila stok menjadi minus (tetap diizinkan).
func recordMessage(base string, data *dto.RecordResponse) string {
	if data.NegativeStock {
		return base + ". Peringatan: stok bahan menjadi minus"
	}
	return base
}

func (h *IngredientHandler) GetMutations(c *gin.Context) {
	req, valid := bindJSON[dto.GetMutationsRequest](c)
	if !valid {
		return
	}
	data, total, err := h.service.GetMutations(&req)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Riwayat stok bahan",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *IngredientHandler) CancelMutation(c *gin.Context) {
	id, valid := bindID(c)
	if !valid {
		return
	}
	req, valid := bindJSON[dto.CancelMutationRequest](c)
	if !valid {
		return
	}
	req.ID = id
	data, err := h.service.CancelMutation(&req, helper.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Entri berhasil dibatalkan", data)
}
