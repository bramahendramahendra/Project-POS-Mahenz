package handler

import (
	dto "pos_api/domain/capital_reconciliation/dto"
	service "pos_api/domain/capital_reconciliation/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type CapitalReconciliationHandler struct {
	service service.CapitalReconciliationServiceInterface
}

func NewCapitalReconciliationHandler(service service.CapitalReconciliationServiceInterface) *CapitalReconciliationHandler {
	return &CapitalReconciliationHandler{service: service}
}

func (h *CapitalReconciliationHandler) GetList(c *gin.Context) {
	req, err := binder.BindJSON[dto.ListRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetList(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Daftar rekonsiliasi modal",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *CapitalReconciliationHandler) GetSummary(c *gin.Context) {
	data, err := h.service.GetSummary()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Ringkasan rekonsiliasi modal",
		Data:    data,
	})
}

func (h *CapitalReconciliationHandler) Resolve(c *gin.Context) {
	req, err := binder.BindJSON[dto.ResolveRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	req.UserID = helper.GetUserID(c)

	if err := h.service.Resolve(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Modal berhasil dikoreksi",
	})
}

func (h *CapitalReconciliationHandler) Skip(c *gin.Context) {
	req, err := binder.BindJSON[dto.SkipRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	req.UserID = helper.GetUserID(c)

	if err := h.service.Skip(&req); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Baris ditandai sudah ditinjau",
	})
}
