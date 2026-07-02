package middleware

import (
	access_service "pos_api/domain/access/service"
	"pos_api/errors"
	log_helper "pos_api/helper/log"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware memastikan role user memiliki akses yang dibutuhkan
// pada menu tertentu berdasarkan data role_menu_access di DB (dengan cache).
//
// Usage:
//
//	PermissionMiddleware(svc, "produk.produk", "can_create")
//
// action: "can_view" | "can_create" | "can_edit" | "can_delete"
func PermissionMiddleware(svc access_service.AccessServiceInterface, menuKey, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			log_helper.SetLog(c, "warn", "Permission Middleware", "user_role tidak ditemukan di context", nil,
				map[string]any{"endpoint": c.Request.RequestURI},
			)
			c.Error(&errors.UnauthenticatedError{Message: "Unauthorized"})
			c.Abort()
			return
		}

		roleName, ok := roleVal.(string)
		if !ok {
			log_helper.SetLog(c, "warn", "Permission Middleware", "user_role bukan tipe string", nil,
				map[string]any{"endpoint": c.Request.RequestURI, "role_value": roleVal},
			)
			c.Error(&errors.UnauthenticatedError{Message: "Unauthorized"})
			c.Abort()
			return
		}

		perm, err := svc.GetPermission(roleName, menuKey)
		if err != nil {
			c.Error(&errors.InternalServerError{Message: "Gagal memuat permission"})
			c.Abort()
			return
		}

		allowed := false
		switch action {
		case "can_view":
			allowed = perm.CanView
		case "can_create":
			allowed = perm.CanCreate
		case "can_edit":
			allowed = perm.CanEdit
		case "can_delete":
			allowed = perm.CanDelete
		}

		if !allowed {
			log_helper.SetLog(c, "warn", "Permission Middleware", "Akses ditolak: permission tidak mencukupi", nil,
				map[string]any{
					"role":     roleName,
					"menu_key": menuKey,
					"action":   action,
					"endpoint": c.Request.RequestURI,
				},
			)
			c.Error(&errors.UnauthorizededError{Message: "Anda tidak memiliki akses ke fitur ini"})
			c.Abort()
			return
		}

		c.Next()
	}
}
