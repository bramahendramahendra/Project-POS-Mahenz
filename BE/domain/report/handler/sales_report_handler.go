package handler

import (
	"net/http"
	"time"

	"pos_api/domain/report/dto"
	"pos_api/domain/report/service"
	global_dto "pos_api/dto"
	"pos_api/helper"
	response_helper "pos_api/helper/response"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	service service.ReportServiceInterface
}

func NewReportHandler(svc service.ReportServiceInterface) *ReportHandler {
	return &ReportHandler{service: svc}
}

func parseFilterParams(c *gin.Context) dto.FilterParams {
	dateFrom := c.DefaultQuery("date_from", time.Now().Format("2006-01-02")+" 00:00:00")
	dateTo := c.DefaultQuery("date_to", time.Now().Format("2006-01-02")+" 23:59:59")
	return dto.FilterParams{DateFrom: dateFrom, DateTo: dateTo}
}

func (h *ReportHandler) GetSalesReport(c *gin.Context) {
	params := parseFilterParams(c)
	data, err := h.service.GetSalesReport(params)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Laporan penjualan",
		Data:    data,
	})
}

func (h *ReportHandler) GetSalesChart(c *gin.Context) {
	params := parseFilterParams(c)
	data, err := h.service.GetSalesChart(params)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Chart penjualan",
		Data:    data,
	})
}

func (h *ReportHandler) ExportSalesReport(c *gin.Context) {
	params := parseFilterParams(c)
	buf, err := h.service.ExportSalesReport(params)
	if err != nil {
		c.Error(err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=laporan-penjualan.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
