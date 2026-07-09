package handler

import (
	"pos_api/domain/pin/dto"
	"pos_api/domain/pin/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"
	"pos_api/pkg/binder"
	validator "pos_api/validation"

	"github.com/gin-gonic/gin"
)

type PinHandler struct {
	service service.PinServiceInterface
}

func NewPinHandler(svc service.PinServiceInterface) *PinHandler {
	return &PinHandler{service: svc}
}

func (h *PinHandler) CheckPin(c *gin.Context) {
	userID := helper.GetUserID(c)

	hasPin, err := h.service.HasPin(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Status PIN",
		Data:    &dto.HasPinResponse{HasPin: hasPin},
	})
}

func (h *PinHandler) SetPin(c *gin.Context) {
	req, err := binder.BindJSON[dto.SetPinRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userID := helper.GetUserID(c)
	if err := h.service.SetPin(userID, req.Pin); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "PIN berhasil diset",
	})
}

func (h *PinHandler) VerifyPin(c *gin.Context) {
	req, err := binder.BindJSON[dto.VerifyPinRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userID := helper.GetUserID(c)
	valid, err := h.service.VerifyPin(userID, req.Pin)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Verifikasi PIN",
		Data:    &dto.VerifyPinResponse{Valid: valid},
	})
}

func (h *PinHandler) ChangePin(c *gin.Context) {
	req, err := binder.BindJSON[dto.ChangePinRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userID := helper.GetUserID(c)
	if err := h.service.ChangePin(userID, req.OldPin, req.NewPin); err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "PIN berhasil diubah",
	})
}
