package handler

import (
	dto "auth_api/domain/app/dto"
	"auth_api/domain/app/service"
	global_dto "auth_api/dto"
	"auth_api/errors"
	"auth_api/helper"
	response_helper "auth_api/helper/response"
	binder "auth_api/pkg/binder"
	validator "auth_api/validation"

	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	service service.AppServiceInterface
}

func NewAppHandler(svc service.AppServiceInterface) *AppHandler {
	return &AppHandler{service: svc}
}

func ok(c *gin.Context, status int, message string, data any) {
	response_helper.WrapResponse(c, status, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// GetMy mengembalikan aplikasi yang boleh dibuka oleh user yang sedang login.
func (h *AppHandler) GetMy(c *gin.Context) {
	data, err := h.service.GetMyApps(helper.GetUserID(c))
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Daftar aplikasi", data)
}

// ---- Khusus admin ----

func (h *AppHandler) GetAll(c *gin.Context) {
	data, err := h.service.GetAll()
	if err != nil {
		c.Error(err)
		return
	}
	ok(c, 200, "Daftar aplikasi (admin)", data)
}

func (h *AppHandler) Create(c *gin.Context) {
	req, err := binder.BindJSON[dto.CreateRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, svcErr := h.service.Create(&req)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}
	ok(c, 201, "Aplikasi dibuat. Simpan kunci aplikasi sekarang, kunci tidak akan ditampilkan lagi", data)
}

func (h *AppHandler) Update(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.IDUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(uriReq); err != nil {
		c.Error(err)
		return
	}

	bodyReq, err := binder.BindJSON[dto.UpdateRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(bodyReq); err != nil {
		c.Error(err)
		return
	}

	if svcErr := h.service.Update(uriReq.ID, &bodyReq); svcErr != nil {
		c.Error(svcErr)
		return
	}
	ok(c, 200, "Aplikasi diperbarui", nil)
}

func (h *AppHandler) ToggleStatus(c *gin.Context) {
	req, err := binder.BindURI[dto.IDUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	if svcErr := h.service.ToggleStatus(req.ID); svcErr != nil {
		c.Error(svcErr)
		return
	}
	ok(c, 200, "Status aplikasi diubah", nil)
}

func (h *AppHandler) RotateKey(c *gin.Context) {
	req, err := binder.BindURI[dto.IDUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, svcErr := h.service.RotateKey(req.ID)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}
	ok(c, 200, "Kunci baru dibuat. Simpan sekarang, kunci lama sudah tidak berlaku", data)
}

func (h *AppHandler) GetUserAccess(c *gin.Context) {
	req, err := binder.BindURI[dto.UserUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, svcErr := h.service.GetUserAccess(req.UserID)
	if svcErr != nil {
		c.Error(svcErr)
		return
	}
	ok(c, 200, "Akses aplikasi user", data)
}

func (h *AppHandler) SetUserAccess(c *gin.Context) {
	uriReq, err := binder.BindURI[dto.UserUriRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(uriReq); err != nil {
		c.Error(err)
		return
	}

	bodyReq, err := binder.BindJSON[dto.SetUserAccessRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(bodyReq); err != nil {
		c.Error(err)
		return
	}

	if svcErr := h.service.SetUserAccess(uriReq.UserID, &bodyReq); svcErr != nil {
		c.Error(svcErr)
		return
	}
	ok(c, 200, "Akses aplikasi user diperbarui", nil)
}
