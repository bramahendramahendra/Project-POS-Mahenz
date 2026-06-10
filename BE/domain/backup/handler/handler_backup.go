package handler

import (
	"path/filepath"
	"strings"

	"pos_api/domain/backup/service"
	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"

	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	service service.BackupServiceInterface
}

func NewBackupHandler(svc service.BackupServiceInterface) *BackupHandler {
	return &BackupHandler{service: svc}
}

// POST /api/backup
func (h *BackupHandler) Create(c *gin.Context) {
	info, err := h.service.CreateBackup()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 201, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Backup berhasil dibuat",
		Data:    info,
	})
}

// GET /api/backup/list
func (h *BackupHandler) GetList(c *gin.Context) {
	result, err := h.service.GetList()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Daftar backup",
		Data:    result,
	})
}

// GET /api/backup/download/:filename
func (h *BackupHandler) Download(c *gin.Context) {
	filename := c.Param("filename")

	// Cegah path traversal
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") || strings.Contains(filename, "..") {
		c.Error(&errors.BadRequestError{Message: "Nama file tidak valid"})
		return
	}
	if !strings.HasSuffix(filename, ".sql") {
		c.Error(&errors.BadRequestError{Message: "Hanya file .sql yang dapat diunduh"})
		return
	}

	filePath := filepath.Join("backups", filepath.Base(filename))
	c.FileAttachment(filePath, filename)
}

// POST /api/restore
func (h *BackupHandler) Restore(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(&errors.BadRequestError{Message: "File SQL wajib disertakan"})
		return
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".sql") {
		c.Error(&errors.BadRequestError{Message: "Hanya file .sql yang diizinkan"})
		return
	}

	if restoreErr := h.service.RestoreBackup(file); restoreErr != nil {
		c.Error(restoreErr)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &global_dto.ResponseParams{
		Code:    helper.StatusOk,
		Status:  true,
		Message: "Restore berhasil dilakukan",
	})
}
