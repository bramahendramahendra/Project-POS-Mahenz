package middleware

import (
	"pos_api/errors"
	log_helper "pos_api/helper/log"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware restricts access to the specified roles.
// Usage: RoleMiddleware("owner", "admin")
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			log_helper.SetLog(c, "warn", "Role Middleware", "user_role tidak ditemukan di context", nil,
				map[string]any{"endpoint": c.Request.RequestURI},
			)
			c.Error(&errors.UnauthenticatedError{Message: "Unauthorized"})
			c.Abort()
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			log_helper.SetLog(c, "warn", "Role Middleware", "user_role bukan tipe string", nil,
				map[string]any{"endpoint": c.Request.RequestURI, "role_value": roleVal},
			)
			c.Error(&errors.UnauthenticatedError{Message: "Unauthorized"})
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		log_helper.SetLog(c, "warn", "Role Middleware", "Akses ditolak: role tidak diizinkan", nil,
			map[string]any{
				"user_role":     userRole,
				"allowed_roles": allowedRoles,
				"endpoint":      c.Request.RequestURI,
			},
		)
		c.Error(&errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke fitur ini"})
		c.Abort()
	}
}
