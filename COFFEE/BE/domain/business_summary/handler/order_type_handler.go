package handler

import (
	dto "coffee_api/domain/business_summary/dto"
	global_dto "coffee_api/dto"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"
	time_helper "coffee_api/helper/time"

	"github.com/gin-gonic/gin"
)

// GetOrderTypes: penjualan per jenis pesanan (Coffee, keputusan 117). Rentang bawaan sama dengan metode bayar.
func (h *BusinessSummaryHandler) GetOrderTypes(c *gin.Context) {
	now := time_helper.GetTimeNow()
	filter := dto.DateRangeFilter{
		StartDate: c.DefaultQuery("start_date", now.AddDate(0, -1, 0).Format("2006-01-02")),
		EndDate:   c.DefaultQuery("end_date", now.Format("2006-01-02")),
	}
	result, err := h.service.GetOrderTypes(filter)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{Code: helper.StatusOk, Status: true, Message: "Jenis pesanan", Data: result})
}
