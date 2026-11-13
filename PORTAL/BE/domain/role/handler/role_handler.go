package handler

import (
	service "portal_api/domain/role/service"
	"portal_api/helper"
	"portal_api/helper/forward"
	"portal_api/pkg/authclient"

	"github.com/gin-gonic/gin"
)

// RoleHandler: role akun (admin, user) dikelola BE User. Portal hanya meneruskan daftar
// pilihan role untuk form user.
type RoleHandler struct {
	service service.RoleServiceInterface
}

func NewRoleHandler(service service.RoleServiceInterface) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) GetOptions(c *gin.Context) {
	forward.Handle(c, func(body []byte) (*authclient.Response, error) {
		return h.service.Forward(c.Request.Context(), helper.GetAuthUserID(c), "/roles/options", body)
	})
}
