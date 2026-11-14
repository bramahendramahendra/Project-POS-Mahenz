package handler

import (
	dto "coffee_api/domain/report/dto"
	global_dto "coffee_api/dto"
	"coffee_api/errors"
	"coffee_api/helper"
	response_helper "coffee_api/helper/response"
	binder "coffee_api/pkg/binder"

	"github.com/gin-gonic/gin"
)

// Laporan Coffee (tahap 4.10): rincian penjualan, bahan baku, barang terbuang.

func bindCoffeeReport(c *gin.Context) (*dto.CoffeeReportRequest, bool) {
	req, err := binder.BindJSON[dto.CoffeeReportRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return nil, false
	}
	return &req, true
}

func respondReport(c *gin.Context, message string, data any, err error) {
	if err != nil {
		c.Error(err)
		return
	}
	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{Code: helper.StatusOk, Status: true, Message: message, Data: data})
}

func (h *ReportHandler) GetSalesBreakdown(c *gin.Context) {
	if req, ok := bindCoffeeReport(c); ok {
		data, err := h.service.GetSalesBreakdown(req)
		respondReport(c, "Rincian penjualan", data, err)
	}
}

func (h *ReportHandler) GetIngredientReport(c *gin.Context) {
	if req, ok := bindCoffeeReport(c); ok {
		data, err := h.service.GetIngredientReport(req)
		respondReport(c, "Laporan bahan baku", data, err)
	}
}

func (h *ReportHandler) GetWasteReport(c *gin.Context) {
	if req, ok := bindCoffeeReport(c); ok {
		data, err := h.service.GetWasteReport(req)
		respondReport(c, "Laporan barang terbuang", data, err)
	}
}
