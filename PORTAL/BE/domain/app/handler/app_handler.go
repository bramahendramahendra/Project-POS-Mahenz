package handler

import (
	"fmt"

	service "portal_api/domain/app/service"
	"portal_api/helper"
	"portal_api/helper/forward"
	"portal_api/pkg/authclient"

	"github.com/gin-gonic/gin"
)

// AppHandler: daftar aplikasi, kunci aplikasi, dan akses user dikelola BE User.
// Portal meneruskan permintaan atas nama user yang sedang login (admin diperiksa BE User).
type AppHandler struct {
	service service.AppServiceInterface
}

func NewAppHandler(svc service.AppServiceInterface) *AppHandler {
	return &AppHandler{service: svc}
}

func (h *AppHandler) send(c *gin.Context, path string, body []byte) (*authclient.Response, error) {
	return h.service.Forward(c.Request.Context(), helper.GetAuthUserID(c), path, body)
}

// GetMy: kartu aplikasi milik user yang sedang login (semua user).
func (h *AppHandler) GetMy(c *gin.Context) {
	forward.Handle(c, func(body []byte) (*authclient.Response, error) { return h.send(c, "/apps/my", body) })
}

func (h *AppHandler) GetAll(c *gin.Context) {
	forward.Handle(c, func(body []byte) (*authclient.Response, error) { return h.send(c, "/apps/list", body) })
}

func (h *AppHandler) Create(c *gin.Context) {
	forward.Handle(c, func(body []byte) (*authclient.Response, error) { return h.send(c, "/apps/create", body) })
}

func (h *AppHandler) Update(c *gin.Context) {
	forward.HandleWithID(c, "id", func(id int, body []byte) (*authclient.Response, error) {
		return h.send(c, fmt.Sprintf("/apps/update/%d", id), body)
	})
}

func (h *AppHandler) ToggleStatus(c *gin.Context) {
	forward.HandleWithID(c, "id", func(id int, body []byte) (*authclient.Response, error) {
		return h.send(c, fmt.Sprintf("/apps/toggle-status/%d", id), body)
	})
}

func (h *AppHandler) RotateKey(c *gin.Context) {
	forward.HandleWithID(c, "id", func(id int, body []byte) (*authclient.Response, error) {
		return h.send(c, fmt.Sprintf("/apps/rotate-key/%d", id), body)
	})
}

func (h *AppHandler) GetUserAccess(c *gin.Context) {
	forward.HandleWithID(c, "user_id", func(id int, body []byte) (*authclient.Response, error) {
		return h.send(c, fmt.Sprintf("/apps/access/%d", id), body)
	})
}

func (h *AppHandler) SetUserAccess(c *gin.Context) {
	forward.HandleWithID(c, "user_id", func(id int, body []byte) (*authclient.Response, error) {
		return h.send(c, fmt.Sprintf("/apps/access/%d/set", id), body)
	})
}
