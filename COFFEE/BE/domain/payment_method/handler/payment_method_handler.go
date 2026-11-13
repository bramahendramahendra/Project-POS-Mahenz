package handler

import (
	service "coffee_api/domain/payment_method/service"
	global_dto "coffee_api/dto"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"

	"github.com/gin-gonic/gin"
)

type PaymentMethodHandler struct {
	service service.PaymentMethodServiceInterface
}

func NewPaymentMethodHandler(service service.PaymentMethodServiceInterface) *PaymentMethodHandler {
	return &PaymentMethodHandler{service: service}
}

func (h *PaymentMethodHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll()
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar metode pembayaran",
		Data:    data,
	})
}
