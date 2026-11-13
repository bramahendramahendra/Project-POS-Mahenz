package routes

import (
	"github.com/gin-gonic/gin"
)

// publicRoutes: tanpa kunci aplikasi. Hanya health check; BE User sendiri hanya
// mendengarkan 127.0.0.1 sehingga tidak terjangkau dari luar server (keputusan 63).
func publicRoutes(r *gin.RouterGroup) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code":    "00",
			"status":  true,
			"message": "OK",
			"data":    nil,
		})
	})
}
