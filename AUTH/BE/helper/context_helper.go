package helper

import "github.com/gin-gonic/gin"

// GetUserID: user yang bertindak (diisi middleware ActingUser dari header X-Acting-User-Id).
func GetUserID(c *gin.Context) int {
	userID, _ := c.Get("user_id")
	return userID.(int)
}

// CallerApp: aplikasi yang memanggil BE User, dikenali dari kunci aplikasi (middleware AppKey).
type CallerApp struct {
	ID   int
	Code string
}

// PortalAppCode: kode aplikasi Portal. Hanya Portal yang boleh mengelola akun dan meminta kode SSO.
const PortalAppCode = "portal"

func GetCallerApp(c *gin.Context) CallerApp {
	app, _ := c.Get("caller_app")
	return app.(CallerApp)
}
