package handler

import (
	"net/http"

	"pos_api/domain/report/dto"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	request_helper "pos_api/helper/request"
	response_helper "pos_api/helper/response"
	binder "pos_api/pkg/binder"

	"github.com/gin-gonic/gin"
)

func (h *ReportHandler) GetCashierReport(c *gin.Context) {
	params := parseFilterParams(c)
	data, err := h.service.GetCashierReport(params)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Laporan kasir",
		Data:    data,
	})
}

func (h *ReportHandler) GetCashierList(c *gin.Context) {
	req, err := binder.BindJSON[dto.CashierReportRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	data, total, err := h.service.GetCashierList(&req)
	if err != nil {
		c.Error(err)
		return
	}
	page, limit, _ := request_helper.NormalizePagination(req.Page, req.Limit, 10, 0)
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:       helper.StatusOk,
		Status:     true,
		Message:    "Daftar kinerja kasir",
		Data:       data,
		Pagination: response_helper.SetPagination(&global_dto.FilterRequestParams{Page: page, Limit: limit}, total),
	})
}

func (h *ReportHandler) ExportCashierReport(c *gin.Context) {
	params := parseFilterParams(c)
	buf, err := h.service.ExportCashierReport(params)
	if err != nil {
		c.Error(err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=laporan-kasir.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
