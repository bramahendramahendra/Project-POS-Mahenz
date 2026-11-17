package handler

import (
	"strconv"

	dto "pos_api/domain/business_summary/dto"
	"pos_api/domain/business_summary/service"
	global_dto "pos_api/dto"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	time_helper "pos_api/helper/time"

	"github.com/gin-gonic/gin"
)

type BusinessSummaryHandler struct {
	service service.BusinessSummaryServiceInterface
}

func NewBusinessSummaryHandler(svc service.BusinessSummaryServiceInterface) *BusinessSummaryHandler {
	return &BusinessSummaryHandler{service: svc}
}

// GET /api/reports/business-summary/stats
func (h *BusinessSummaryHandler) GetStats(c *gin.Context) {
	period := c.DefaultQuery("period", "today")
	result, err := h.service.GetStats(period)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Statistik ringkasan bisnis",
		Data:    result,
	})
}

// GET /api/reports/business-summary/sales-trend
func (h *BusinessSummaryHandler) GetSalesTrend(c *gin.Context) {
	period := c.DefaultQuery("period", "7days")
	result, err := h.service.GetSalesTrend(period)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Tren penjualan",
		Data:    result,
	})
}

// GET /api/reports/business-summary/top-products
func (h *BusinessSummaryHandler) GetTopProducts(c *gin.Context) {
	now := time_helper.GetTimeNow()
	startDate := c.DefaultQuery("start_date", now.AddDate(0, -1, 0).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", now.Format("2006-01-02"))
	sortBy := c.DefaultQuery("sort_by", "quantity")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	}

	filter := dto.DateRangeFilter{
		StartDate: startDate,
		EndDate:   endDate,
		SortBy:    sortBy,
		Limit:     limit,
	}

	result, err := h.service.GetTopProducts(filter)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Produk terlaris",
		Data:    result,
	})
}
