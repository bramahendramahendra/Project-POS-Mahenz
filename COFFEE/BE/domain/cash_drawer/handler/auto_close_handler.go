package handler

import (
	global_dto "coffee_api/dto"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"

	"github.com/gin-gonic/gin"
)

// AutoCloseYesterday menjalankan penutupan kas hari sebelumnya secara manual (sama persis dengan
// scheduler tengah malam). Kas hari yang masih punya pesanan terbuka dibiarkan tertunda (keputusan 111).
// Dipakai admin/owner bila server mati saat tengah malam, dan untuk uji di dev (tahap 4.9).
func (h *CashDrawerHandler) AutoCloseYesterday(c *gin.Context) {
	closed, err := h.service.AutoCloseYesterday()
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Penutupan kas hari sebelumnya selesai",
		Data:    map[string]int{"closed": closed},
	})
}
