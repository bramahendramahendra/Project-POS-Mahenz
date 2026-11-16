package helper

import "github.com/gin-gonic/gin"

// GetUserID: id profil lokal Portal (tabel users) dari sesi yang sedang login.
func GetUserID(c *gin.Context) int {
	userID, _ := c.Get("user_id")
	return userID.(int)
}

// GetAuthUserID: id akun di BE User dari sesi yang sedang login. Dikirim ke BE User
// sebagai X-Acting-User-Id.
func GetAuthUserID(c *gin.Context) int {
	authUserID, _ := c.Get("auth_user_id")
	return authUserID.(int)
}
