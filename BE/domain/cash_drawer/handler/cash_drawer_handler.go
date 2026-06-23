package handler

import (
	dto "pos_api/domain/cash_drawer/dto"
	service "pos_api/domain/cash_drawer/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type CashDrawerHandler struct {
	service service.CashDrawerServiceInterface
}

func NewCashDrawerHandler(service service.CashDrawerServiceInterface) *CashDrawerHandler {
	return &CashDrawerHandler{service: service}
}

func (h *CashDrawerHandler) GetCurrent(c *gin.Context) {
	userID := helper.GetUserID(c)

	data, err := h.service.GetCurrent(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Status kas harian",
		Data:    data,
	})
}

func (h *CashDrawerHandler) GetMyCash(c *gin.Context) {
	userID := helper.GetUserID(c)

	data, err := h.service.GetMyCash(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Kas saya",
		Data:    data,
	})
}

func (h *CashDrawerHandler) GetHistory(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetHistoryRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	role := helper.GetUserRole(c)
	requestingUserID := helper.GetUserID(c)

	if role != "owner" && role != "admin" {
		req.UserID = &requestingUserID
	}

	data, total, err := h.service.GetHistory(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Riwayat kas",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: req.Page, Limit: req.Limit}, total),
	})
}

func (h *CashDrawerHandler) GetByID(c *gin.Context) {
	req, err := binder.BindURI[dto.GetByIDRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetByID(req.ID, helper.GetUserID(c), helper.GetUserRole(c))
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Detail kas",
		Data:    data,
	})
}

func (h *CashDrawerHandler) Open(c *gin.Context) {
	req, err := binder.BindJSON[dto.OpenRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userID := helper.GetUserID(c)

	data, err := h.service.Open(userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Kas berhasil dibuka",
		Data:    data,
	})
}

func (h *CashDrawerHandler) Close(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.CloseUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.CloseRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.ID

	data, err := h.service.Close(req.ID, &req, helper.GetUserID(c), helper.GetUserRole(c))
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Kas berhasil ditutup",
		Data:    data,
	})
}

func (h *CashDrawerHandler) UpdateSales(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.UpdateSalesUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.UpdateSalesRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.ID

	if err := h.service.UpdateSales(req.ID, &req, helper.GetUserID(c), helper.GetUserRole(c)); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Data penjualan berhasil diperbarui",
	})
}

func (h *CashDrawerHandler) UpdateExpenses(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.UpdateExpensesUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.UpdateExpensesRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.ID = uriReq.ID

	if err := h.service.UpdateExpenses(req.ID, &req, helper.GetUserID(c), helper.GetUserRole(c)); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Data pengeluaran berhasil diperbarui",
	})
}
