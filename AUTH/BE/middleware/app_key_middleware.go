package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"auth_api/errors"
	"auth_api/helper"

	"github.com/gin-gonic/gin"
)

// AppLookup: cari aplikasi aktif dari hash kunci (dipenuhi repo aplikasi).
type AppLookup interface {
	GetActiveByKeyHash(keyHash string) (id int, code string, found bool, err error)
}

// UserLookup: data user yang bertindak (dipenuhi repo auth).
type UserLookup interface {
	GetActingUser(id int) (role string, isActive, mustChangePassword, found bool, err error)
}

// Pesan seragam agar pemanggil tanpa kunci yang benar tidak mendapat petunjuk apa pun.
const invalidAppKey = "Kunci aplikasi tidak valid"

// AppKeyMiddleware: setiap permintaan ke API internal wajib membawa kunci aplikasi (X-App-Key).
// Aplikasi pemanggil disimpan di konteks (helper.GetCallerApp).
func AppKeyMiddleware(apps AppLookup) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := strings.TrimSpace(c.GetHeader("X-App-Key"))
		if key == "" {
			c.Error(&errors.UnauthenticatedError{Message: invalidAppKey})
			c.Abort()
			return
		}
		sum := sha256.Sum256([]byte(key))
		id, code, found, err := apps.GetActiveByKeyHash(hex.EncodeToString(sum[:]))
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		if !found {
			c.Error(&errors.UnauthenticatedError{Message: invalidAppKey})
			c.Abort()
			return
		}
		c.Set("caller_app", helper.CallerApp{ID: id, Code: code})
		c.Next()
	}
}

// PortalOnlyMiddleware: endpoint pengelolaan akun dan peluncuran SSO hanya untuk Portal.
func PortalOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if helper.GetCallerApp(c).Code != helper.PortalAppCode {
			c.Error(&errors.UnauthorizededError{Message: "Endpoint ini hanya untuk Portal"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ActingUserMiddleware: Portal menyebut user yang bertindak lewat header X-Acting-User-Id
// (setelah memeriksa sesi Portal-nya sendiri). requireAdmin: hanya role admin, dan
// akun yang wajib ganti password tidak boleh memakai fitur admin.
func ActingUserMiddleware(users UserLookup, requireAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(strings.TrimSpace(c.GetHeader("X-Acting-User-Id")))
		if err != nil || id <= 0 {
			c.Error(&errors.UnauthenticatedError{Message: "User yang bertindak tidak dikenali"})
			c.Abort()
			return
		}
		role, isActive, mustChange, found, err := users.GetActingUser(id)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		if !found || !isActive {
			c.Error(&errors.UnauthenticatedError{Message: "Akun tidak aktif"})
			c.Abort()
			return
		}
		if requireAdmin {
			if role != "admin" {
				c.Error(&errors.UnauthorizededError{Message: "Akses ditolak"})
				c.Abort()
				return
			}
			if mustChange {
				c.Error(&errors.UnauthorizededError{Message: "Anda wajib mengganti password terlebih dahulu"})
				c.Abort()
				return
			}
		}
		c.Set("user_id", id)
		c.Set("user_role", role)
		c.Next()
	}
}
