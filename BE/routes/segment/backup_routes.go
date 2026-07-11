package segment

import (
	backup_handler "pos_api/domain/backup/handler"
	backup_service "pos_api/domain/backup/service"
	middleware "pos_api/middleware"

	"github.com/gin-gonic/gin"
)

func BackupRoutes(r *gin.RouterGroup) {
	backupService := backup_service.NewBackupService()
	backupHandler := backup_handler.NewBackupHandler(backupService)

	svc := newAccessService()
	perm := func(action string) gin.HandlerFunc {
		return middleware.PermissionMiddleware(svc, "sistem.backup", action)
	}

	g := r.Group("/backup")
	{
		g.POST("", perm("can_create"), backupHandler.Create)
		g.GET("/list", perm("can_view"), backupHandler.GetList)
		g.GET("/download/:filename", perm("can_view"), backupHandler.Download)
	}

	// Restore menimpa seluruh data yang ada — dipetakan ke can_delete, slot permission
	// yang di aplikasi ini sudah jadi konvensi untuk aksi paling sensitif/destruktif.
	r.POST("/restore", perm("can_delete"), backupHandler.Restore)
}
