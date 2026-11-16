package handler

import (
	service "auth_api/domain/role/service"
	global_dto "auth_api/dto"
	"auth_api/helper"
	response_helper "auth_api/helper/response"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service service.RoleServiceInterface
}

func NewRoleHandler(service service.RoleServiceInterface) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) GetOptions(c *gin.Context) {
	data, err := h.service.GetActiveOptions()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Opsi role aktif",
		Data:    data,
	})
}
