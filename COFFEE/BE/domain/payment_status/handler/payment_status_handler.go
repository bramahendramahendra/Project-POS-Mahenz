package handler

import (
	service "coffee_api/domain/payment_status/service"
	global_dto "coffee_api/dto"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"

	"github.com/gin-gonic/gin"
)

type PaymentStatusHandler struct {
	service service.PaymentStatusServiceInterface
}

func NewPaymentStatusHandler(service service.PaymentStatusServiceInterface) *PaymentStatusHandler {
	return &PaymentStatusHandler{service: service}
}

func (h *PaymentStatusHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll()
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar status pembayaran",
		Data:    data,
	})
}
