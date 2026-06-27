package handler

import (
	"net/http"

	"pos_api/domain/report/dto"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"

	"github.com/gin-gonic/gin"
)

func (h *ReportHandler) GetStockReport(c *gin.Context) {
	data, err := h.service.GetStockReport()
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Laporan stok",
		Data:    data,
	})
}

func (h *ReportHandler) GetStockList(c *gin.Context) {
	req, err := binder.BindJSON[dto.StockListRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	data, total, err := h.service.GetStockList(&req)
	if err != nil {
		c.Error(err)
		return
	}
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Daftar stok produk",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: page, Limit: limit}, total),
	})
}

func (h *ReportHandler) GetStockSummaryData(c *gin.Context) {
	req, err := binder.BindJSON[dto.StockSummaryRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	data, err := h.service.GetStockSummaryData(&req)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Ringkasan stok",
		Data:    data,
	})
}

func (h *ReportHandler) ExportStockReport(c *gin.Context) {
	buf, err := h.service.ExportStockReport()
	if err != nil {
		c.Error(err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=laporan-stok.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
