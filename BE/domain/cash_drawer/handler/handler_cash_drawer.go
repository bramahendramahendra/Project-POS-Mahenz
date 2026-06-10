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

	res, err := h.service.GetCurrent(userID)
	if err != nil {
		c.Error(err)
		return
	}

	msg := "Success"
	if res == nil {
		msg = "Tidak ada kas yang terbuka"
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: msg,
		Data:    res,
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

	res, err := h.service.GetByID(req.ID, helper.GetUserID(c), helper.GetUserRole(c))
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Detail kas",
		Data:    res,
	})
}

func (h *CashDrawerHandler) Open(c *gin.Context) {
	req, err := binder.BindJSON[dto.OpenRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userID := helper.GetUserID(c)

	res, svcErr := h.service.Open(userID, &req)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Kas berhasil dibuka",
		Data:    res,
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

	res, svcErr := h.service.Close(req.ID, &req, helper.GetUserID(c), helper.GetUserRole(c))
	if svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Kas berhasil ditutup",
		Data:    res,
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

	if svcErr := h.service.UpdateSales(req.ID, &req, helper.GetUserID(c), helper.GetUserRole(c)); svcErr != nil {
		c.Error(svcErr)
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

	if svcErr := h.service.UpdateExpenses(req.ID, &req, helper.GetUserID(c), helper.GetUserRole(c)); svcErr != nil {
		c.Error(svcErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Data pengeluaran berhasil diperbarui",
	})
}
