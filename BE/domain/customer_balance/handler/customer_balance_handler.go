package handler

import (
	"net/http"

	"pos_api/domain/customer_balance/dto"
	"pos_api/domain/customer_balance/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	"pos_api/pkg/binder"

	"github.com/gin-gonic/gin"
)

type CustomerBalanceHandler struct {
	service service.CustomerBalanceServiceInterface
}

func NewCustomerBalanceHandler(service service.CustomerBalanceServiceInterface) *CustomerBalanceHandler {
	return &CustomerBalanceHandler{service: service}
}

type balanceURI struct {
	ID int `uri:"id" validate:"required,min=1"`
}

// POST /customers/:id/balance/topup
func (h *CustomerBalanceHandler) Topup(c *gin.Context) {
	uri, err := binder.BindURI[balanceURI](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.TopupRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.CustomerID = uri.ID
	req.UserID = c.GetInt("user_id")

	data, err := h.service.Topup(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, http.StatusOK, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Top-up saldo berhasil",
		Data:    data,
	})
}

// POST /customers/:id/balance/refund
func (h *CustomerBalanceHandler) Refund(c *gin.Context) {
	uri, err := binder.BindURI[balanceURI](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.RefundRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.CustomerID = uri.ID
	req.UserID = c.GetInt("user_id")

	data, err := h.service.Refund(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, http.StatusOK, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Tarik saldo berhasil",
		Data:    data,
	})
}

// POST /customers/:id/balance/history
func (h *CustomerBalanceHandler) History(c *gin.Context) {
	uri, err := binder.BindURI[balanceURI](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	req, err := binder.BindJSON[dto.HistoryRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.CustomerID = uri.ID

	data, total, err := h.service.GetHistory(&req)
	if err != nil {
		c.Error(err)
		return
	}

	pagination := response_helper.SetPagination(&global_dto.FilterRequestParams{
		Page:  req.Page,
		Limit: req.Limit,
	}, total)

	response_helper.WrapResponse(c, http.StatusOK, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "OK",
		Data:       data,
		Pagination: pagination,
	})
}

// POST /transactions/save-to-balance
func (h *CustomerBalanceHandler) SaveToBalance(c *gin.Context) {
	req, err := binder.BindJSON[dto.SaveToBalanceRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	req.UserID = c.GetInt("user_id")

	data, err := h.service.SaveToBalance(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, http.StatusOK, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Kembalian berhasil disimpan ke saldo",
		Data:    data,
	})
}
