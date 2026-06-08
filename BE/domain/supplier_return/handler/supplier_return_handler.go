package handler

import (
	"strconv"

	dto_supplier_return "pos_api/domain/supplier_return/dto"
	service_supplier_return "pos_api/domain/supplier_return/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	"pos_api/validation"

	"github.com/gin-gonic/gin"
)

type SupplierReturnHandler struct {
	service service_supplier_return.SupplierReturnService
}

func NewSupplierReturnHandler(service service_supplier_return.SupplierReturnService) *SupplierReturnHandler {
	return &SupplierReturnHandler{service: service}
}

// POST /supplier-returns/list
func (h *SupplierReturnHandler) GetAll(c *gin.Context) {
	req, err := binder.BindJSON[dto_supplier_return.SupplierReturnListRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	filter := &dto_supplier_return.SupplierReturnFilter{
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		SupplierID: req.SupplierID,
		Status:     req.Status,
		Page:       req.Page,
		Limit:      req.Limit,
	}

	items, total, err := h.service.GetAll(filter)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Daftar retur supplier",
		Data:       items,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, int64(total)),
	})
}

func (h *SupplierReturnHandler) GetByID(c *gin.Context) {
	id, err := parseReturnID(c)
	if err != nil {
		c.Error(err)
		return
	}

	item, svcErr := h.service.GetByID(id)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Detail retur supplier",
		Data:    item,
	})
}

// POST /supplier-returns/create
func (h *SupplierReturnHandler) Create(c *gin.Context) {
	req, err := binder.BindJSON[dto_supplier_return.CreateSupplierReturnRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	item, svcErr := h.service.Create(&req, uid)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusCreated,
		Status:  true,
		Message: "Retur supplier berhasil dibuat",
		Data:    item,
	})
}

// POST /supplier-returns/update-status/:id
func (h *SupplierReturnHandler) UpdateStatus(c *gin.Context) {
	id, err := parseReturnID(c)
	if err != nil {
		c.Error(err)
		return
	}

	req, bindErr := binder.BindJSON[dto_supplier_return.UpdateStatusRequest](c)
	if bindErr != nil {
		c.Error(&errors.BadRequestError{Message: bindErr.Error()})
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	if svcErr := h.service.UpdateStatus(id, &req, uid); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Status retur supplier berhasil diperbarui",
	})
}

// POST /supplier-returns/delete/:id
func (h *SupplierReturnHandler) Delete(c *gin.Context) {
	id, err := parseReturnID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if svcErr := h.service.Delete(id); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Retur supplier berhasil dihapus",
	})
}

func parseReturnID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		return 0, &errors.BadRequestError{Message: "ID tidak valid"}
	}
	return id, nil
}
