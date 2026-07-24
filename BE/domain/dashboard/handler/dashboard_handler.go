package handler

import (
	"strconv"

	"pos_api/domain/dashboard/service"
	global_dto "pos_api/dto"
	"pos_api/helper"
	response_helper "pos_api/helper/response"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service service.DashboardServiceInterface
}

func NewDashboardHandler(svc service.DashboardServiceInterface) *DashboardHandler {
	return &DashboardHandler{service: svc}
}

// GET /api/dashboard/recent-transactions?limit=5 — transaksi terakhir milik user yang login.
func (h *DashboardHandler) GetRecentTransactions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if limit <= 0 {
		limit = 5
	}

	result, err := h.service.GetRecentTransactions(uid, limit)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Transaksi terakhir",
		Data:    result,
	})
}

// GET /api/dashboard/today-summary — jumlah & total transaksi milik user yang login, hari ini.
func (h *DashboardHandler) GetTodaySummary(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid, _ := userID.(int)

	result, err := h.service.GetTodaySummary(uid)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Performa hari ini",
		Data:    result,
	})
}
