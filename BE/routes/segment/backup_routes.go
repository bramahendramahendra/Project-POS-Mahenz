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

	g := r.Group("/backup", middleware.RoleMiddleware("owner", "admin"))
	{
		g.POST("", backupHandler.Create)
		g.GET("/list", backupHandler.GetList)
		g.GET("/download/:filename", backupHandler.Download)
	}

	r.POST("/restore", middleware.RoleMiddleware("admin"), backupHandler.Restore)
}
