package handler

import (
	"strconv"

	"pos_api/domain/sync/dto"
	"pos_api/domain/sync/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	"pos_api/pkg/binder"

	"github.com/gin-gonic/gin"
)

type SyncHandler struct {
	service service.SyncServiceInterface
}

func NewSyncHandler(svc service.SyncServiceInterface) *SyncHandler {
	return &SyncHandler{service: svc}
}

func (h *SyncHandler) PushSync(c *gin.Context) {
	req, err := binder.BindJSON[dto.PushSyncRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.PushSync(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Push sync berhasil diproses",
		Data:    data,
	})
}

func (h *SyncHandler) GetConflictCount(c *gin.Context) {
	count, err := h.service.CountPendingConflicts()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Jumlah konflik pending",
		Data:    map[string]int{"count": count},
	})
}

func (h *SyncHandler) GetConflicts(c *gin.Context) {
	filter := &dto.ConflictFilter{
		Status: c.DefaultQuery("status", "pending"),
		Page:   parseIntQuery(c.DefaultQuery("page", "1"), 1),
		Limit:  parseIntQuery(c.DefaultQuery("limit", "20"), 20),
	}

	result, err := h.service.GetConflicts(filter)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar konflik sync",
		Data:    result,
	})
}

func (h *SyncHandler) ResolveConflict(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(&errors.BadRequestError{Message: "ID tidak valid"})
		return
	}

	req, err := binder.BindJSON[dto.ResolveConflictRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userID := helper.GetUserID(c)

	if err := h.service.ResolveConflict(id, userID, req.Action); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Konflik berhasil diselesaikan",
	})
}

func (h *SyncHandler) GetQueue(c *gin.Context) {
	filter := &dto.QueueFilter{
		DeviceID:   c.Query("device_id"),
		Status:     c.Query("status"),
		EntityType: c.Query("entity_type"),
		Page:       parseIntQuery(c.Query("page"), 1),
		Limit:      parseIntQuery(c.Query("limit"), 20),
	}

	result, err := h.service.GetQueue(filter)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Antrian sync",
		Data:    result,
	})
}

func (h *SyncHandler) GetHistory(c *gin.Context) {
	filter := &dto.HistoryFilter{
		DeviceID:  c.Query("device_id"),
		StartDate: c.Query("start_date"),
		EndDate:   c.Query("end_date"),
		Page:      parseIntQuery(c.Query("page"), 1),
		Limit:     parseIntQuery(c.Query("limit"), 20),
	}

	result, err := h.service.GetHistory(filter)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Riwayat sync",
		Data:    result,
	})
}

func parseIntQuery(val string, def int) int {
	if v, err := strconv.Atoi(val); err == nil && v > 0 {
		return v
	}
	return def
}
