package handler

import (
	"fmt"

	service "portal_api/domain/user/service"
	"portal_api/helper"
	"portal_api/helper/forward"
	"portal_api/pkg/authclient"

	"github.com/gin-gonic/gin"
)

// UserHandler: kelola user Portal. Permintaan diteruskan ke BE User atas nama user yang
// sedang login; BE User yang memeriksa role admin, validasi, dan aturan (admin terakhir,
// akun sendiri, dst). Jawaban BE User dikembalikan apa adanya.
type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	forward.Handle(c, func(body []byte) (*authclient.Response, error) {
		return h.service.Forward(c.Request.Context(), helper.GetAuthUserID(c), "/users/list", body)
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	forward.HandleWithID(c, "id", func(id int, body []byte) (*authclient.Response, error) {
		return h.service.Forward(c.Request.Context(), helper.GetAuthUserID(c), fmt.Sprintf("/users/detail/%d", id), body)
	})
}

func (h *UserHandler) Create(c *gin.Context) {
	forward.Handle(c, func(body []byte) (*authclient.Response, error) {
		return h.service.Forward(c.Request.Context(), helper.GetAuthUserID(c), "/users/create", body)
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	forward.HandleWithID(c, "id", func(id int, body []byte) (*authclient.Response, error) {
		return h.service.Forward(c.Request.Context(), helper.GetAuthUserID(c), fmt.Sprintf("/users/update/%d", id), body)
	})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	forward.HandleWithID(c, "id", func(id int, body []byte) (*authclient.Response, error) {
		return h.service.ResetPassword(c.Request.Context(), helper.GetAuthUserID(c), id, body)
	})
}

func (h *UserHandler) Delete(c *gin.Context) {
	forward.HandleWithID(c, "id", func(id int, body []byte) (*authclient.Response, error) {
		return h.service.Delete(c.Request.Context(), helper.GetAuthUserID(c), id, body)
	})
}

func (h *UserHandler) ToggleStatus(c *gin.Context) {
	forward.HandleWithID(c, "id", func(id int, body []byte) (*authclient.Response, error) {
		return h.service.ToggleStatus(c.Request.Context(), helper.GetAuthUserID(c), id, body)
	})
}
