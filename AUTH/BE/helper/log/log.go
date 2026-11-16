package log_helper

import (
	"strings"

	global_dto "auth_api/dto"
	"auth_api/helper"
	request_helper "auth_api/helper/request"
	time_helper "auth_api/helper/time"
	"auth_api/pkg/logger"

	"github.com/gin-gonic/gin"
)

// FromContext membangun LogEntry dari gin.Context.
// Gunakan ini di middleware dan handler yang memiliki akses ke *gin.Context.
func FromContext(c *gin.Context, ctx, scope, message string) global_dto.LogEntry {
	reqData := request_helper.GetRequestContextData(c)
	startTime := reqData["start_time"]
	return global_dto.LogEntry{
		Message:   message,
		Context:   ctx,
		Scope:     scope,
		RequestId: reqData["req_id"],
		Method:    c.Request.Method,
		Endpoint:  c.Request.Host + c.Request.RequestURI,
		StartTime: startTime,
		EndTime:   time_helper.GetEndTime(startTime),
	}
}

// FromBackground membangun LogEntry untuk konteks non-HTTP (service layer, main, scheduler).
// Menggenerate requestId baru dan menghitung waktu secara mandiri.
func FromBackground(ctx, scope, message string) global_dto.LogEntry {
	startTime := time_helper.GetTimeWithFormat()
	return global_dto.LogEntry{
		Message:   message,
		Context:   ctx,
		Scope:     scope,
		RequestId: helper.GenerateUniqueId(),
		StartTime: startTime,
		EndTime:   time_helper.GetEndTime(startTime),
	}
}

// SetLogData adalah fungsi utama untuk logging dari gin.Context.
// logType: "info", "warn", "error", "debug"
func SetLogData(c *gin.Context, logType, ctx, scope, message string, stacktrace any, data any) {
	entry := FromContext(c, ctx, scope, message)
	entry.Data = data
	if st, ok := stacktrace.(string); ok {
		entry.Stacktrace = st
	}
	dispatchLog(logType, entry)
}

// SetLog adalah shorthand SetLogData dengan context "Internal Log".
func SetLog(c *gin.Context, logType, scope, message string, stacktrace, data any) {
	SetLogData(c, logType, "Internal Log", scope, message, stacktrace, data)
}

// LogError logs an already-built LogEntry at Error level.
func LogError(entry global_dto.LogEntry) { logger.Log.Error(entry) }

// LogWarn logs an already-built LogEntry at Warn level.
func LogWarn(entry global_dto.LogEntry) { logger.Log.Warn(entry) }

// LogInfo logs an already-built LogEntry at Info level.
func LogInfo(entry global_dto.LogEntry) { logger.Log.Info(entry) }

// LogDebug logs an already-built LogEntry at Debug level.
func LogDebug(entry global_dto.LogEntry) { logger.Log.Debug(entry) }

func dispatchLog(logType string, entry global_dto.LogEntry) {
	switch strings.ToLower(logType) {
	case "info":
		logger.Log.Info(entry)
	case "warn":
		logger.Log.Warn(entry)
	case "error":
		logger.Log.Error(entry)
	default:
		logger.Log.Debug(entry)
	}
}
