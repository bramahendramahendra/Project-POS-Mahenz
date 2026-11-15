package request_helper

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

func ReadRequestBody(c *gin.Context) (string, error) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", err
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes), nil
}

func GetRequestContextData(c *gin.Context) map[string]string {
	return map[string]string{
		"start_time":     c.GetString("start_time"),
		"req_id":         c.GetString("req_id"),
		"req_header_str": c.GetString("req_header_str"),
		"req_body_str":   c.GetString("req_body_str"),
	}
}
