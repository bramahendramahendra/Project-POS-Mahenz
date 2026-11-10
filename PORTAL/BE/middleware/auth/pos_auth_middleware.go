package middleware

import (
	"strings"

	service_auth "portal_api/domain/auth/service"
	"portal_api/errors"

	"github.com/gin-gonic/gin"
)

// Endpoint yang tetap boleh diakses saat user wajib mengganti password.
var allowedWhenMustChangePassword = map[string]bool{
	"/api/auth/me":              true,
	"/api/auth/logout":          true,
	"/api/auth/change-password": true,
}

// POSBearerAuthMiddleware validates JWT + single active session via sessions table.
// Bila user wajib mengganti password (password sementara hasil reset admin),
// hanya endpoint pada allowedWhenMustChangePassword yang boleh diakses.
func POSBearerAuthMiddleware(authService service_auth.AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Error(&errors.UnauthenticatedError{Message: "Token tidak ditemukan"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		session, err := authService.ValidateToken(token)
		if err != nil {
			c.Error(&errors.UnauthenticatedError{Message: err.Error()})
			c.Abort()
			return
		}

		if session.MustChangePassword && !allowedWhenMustChangePassword[c.FullPath()] {
			c.Error(&errors.UnauthorizededError{Message: "Anda wajib mengganti password terlebih dahulu"})
			c.Abort()
			return
		}

		c.Set("user_id", session.UserID)
		c.Set("user_role", session.UserRole)
		c.Set("device_info", session.DeviceInfo)
		c.Set("token", token)

		c.Next()
	}
}
