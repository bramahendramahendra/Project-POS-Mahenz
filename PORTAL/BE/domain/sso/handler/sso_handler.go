package handler

import (
	service "portal_api/domain/sso/service"
	"portal_api/helper"
	"portal_api/helper/forward"
	"portal_api/pkg/authclient"

	"github.com/gin-gonic/gin"
)

// SSOHandler: kode SSO dibuat oleh BE User. Portal hanya meneruskan permintaan launch atas
// nama user yang sedang login. Penukaran kode (exchange) dilakukan aplikasi tujuan langsung
// ke BE User, bukan lewat Portal.
type SSOHandler struct {
	service service.SSOServiceInterface
}

func NewSSOHandler(svc service.SSOServiceInterface) *SSOHandler {
	return &SSOHandler{service: svc}
}

func (h *SSOHandler) Launch(c *gin.Context) {
	forward.Handle(c, func(body []byte) (*authclient.Response, error) {
		return h.service.Forward(c.Request.Context(), helper.GetAuthUserID(c), "/sso/launch", body)
	})
}
