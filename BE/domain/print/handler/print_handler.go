package handler

import (
	"pos_api/domain/print/dto"
	"pos_api/domain/print/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"

	"github.com/gin-gonic/gin"
)

type PrintHandler struct {
	service service.PrintServiceInterface
}

func NewPrintHandler(svc service.PrintServiceInterface) *PrintHandler {
	return &PrintHandler{service: svc}
}

// POST /api/print/receipt
// Mencetak struk langsung ke printer thermal via COM port (ESC/POS).
func (h *PrintHandler) PrintReceipt(c *gin.Context) {
	req, err := binder.BindJSON[dto.PrintReceiptRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := h.service.PrintReceipt(&req); err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Struk berhasil dicetak",
	})
}
