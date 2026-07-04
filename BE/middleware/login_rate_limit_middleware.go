package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	global_dto "pos_api/dto"
	"pos_api/errors"
	"pos_api/helper"
	response_helper "pos_api/helper/response"

	"github.com/gin-gonic/gin"
)

const (
	loginRateLimitMaxAttempts = 5
	loginRateLimitWindow      = 15 * time.Minute
)

type loginAttemptTracker struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
}

var loginTracker = &loginAttemptTracker{attempts: make(map[string][]time.Time)}

// isBlocked membuang percobaan yang sudah di luar window, lalu mengembalikan true jika
// jumlah percobaan gagal yang masih berlaku sudah mencapai batas maksimum.
func (t *loginAttemptTracker) isBlocked(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	valid := t.attempts[key][:0]
	for _, ts := range t.attempts[key] {
		if now.Sub(ts) < loginRateLimitWindow {
			valid = append(valid, ts)
		}
	}
	t.attempts[key] = valid
	return len(valid) >= loginRateLimitMaxAttempts
}

func (t *loginAttemptTracker) recordFailure(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.attempts[key] = append(t.attempts[key], time.Now())
}

func (t *loginAttemptTracker) reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.attempts, key)
}

// LoginRateLimitMiddleware membatasi percobaan login gagal per kombinasi IP+username
// dalam sebuah window waktu, untuk mengurangi risiko brute-force password.
func LoginRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var body struct {
			Username string `json:"username"`
		}
		_ = json.Unmarshal(bodyBytes, &body)

		key := c.ClientIP() + "|" + strings.ToLower(strings.TrimSpace(body.Username))

		if loginTracker.isBlocked(key) {
			response_helper.WrapResponse(c, http.StatusTooManyRequests, "json", &global_dto.ResponseParams{
				Code:    helper.StatusTooManyRequests,
				Status:  false,
				Message: "Terlalu banyak percobaan login gagal. Coba lagi beberapa saat lagi.",
			})
			c.Abort()
			return
		}

		c.Next()

		if len(c.Errors) > 0 {
			if _, ok := c.Errors.Last().Err.(*errors.UnauthenticatedError); ok {
				loginTracker.recordFailure(key)
			}
			return
		}

		loginTracker.reset(key)
	}
}
