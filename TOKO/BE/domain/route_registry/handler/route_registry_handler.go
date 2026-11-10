package handler

import (
	service "pos_api/domain/route_registry/service"
	global_dto "pos_api/dto"
	"pos_api/helper"
	response_helper "pos_api/helper/response"

	"github.com/gin-gonic/gin"
)

type RouteRegistryHandler struct {
	service service.RouteRegistryServiceInterface
}

func NewRouteRegistryHandler(service service.RouteRegistryServiceInterface) *RouteRegistryHandler {
	return &RouteRegistryHandler{service: service}
}

func (h *RouteRegistryHandler) GetOptions(c *gin.Context) {
	data, err := h.service.GetOptions()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Opsi path route",
		Data:    data,
	})
}
