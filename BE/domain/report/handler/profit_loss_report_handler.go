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

func (h *ReportHandler) GetProfitLoss(c *gin.Context) {
	params := parseFilterParams(c)
	data, err := h.service.GetProfitLoss(params)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Laporan laba rugi",
		Data:    data,
	})
}

func (h *ReportHandler) GetProfitLossData(c *gin.Context) {
	req, err := binder.BindJSON[dto.ProfitLossRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	data, err := h.service.GetProfitLossData(&req)
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Laporan laba rugi",
		Data:    data,
	})
}

func (h *ReportHandler) ExportProfitLoss(c *gin.Context) {
	params := parseFilterParams(c)
	buf, err := h.service.ExportProfitLoss(params)
	if err != nil {
		c.Error(err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=laporan-laba-rugi.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
