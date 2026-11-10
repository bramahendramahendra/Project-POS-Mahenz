package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"pos_api/config"
	global_dto "pos_api/dto"
	time_helper "pos_api/helper/time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log        *Logger
	timeFormat string = "02012006"
)

type Logger struct {
	*zap.Logger
	mu   sync.RWMutex
	date string
}

func New() *Logger {
	currDate := time_helper.GetTimeNow().Format(timeFormat)
	cfg := config.Cfg.Log
	minLevel := parseLevel(cfg.Level)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "message",
		CallerKey:    "caller",
		NameKey:      "logger",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		LineEnding:   zapcore.DefaultLineEnding,
	}

	errorWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("./%s/error/error_%s.log", cfg.Path, currDate),
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   true,
	})
	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		errorWriter,
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= zapcore.ErrorLevel
		}),
	)

	appWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("./%s/app/app_%s.log", cfg.Path, currDate),
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   true,
	})
	appCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		appWriter,
		zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= minLevel && l < zapcore.ErrorLevel
		}),
	)

	core := zapcore.NewTee(errorCore, appCore)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &Logger{Logger: zapLogger, date: currDate}
}

// StartRotationWatcher harus dijalankan sebagai goroutine saat startup.
// Mengecek pergantian tanggal setiap menit dan rotate file log jika perlu.
func StartRotationWatcher() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		Log.rotateIfNeeded()
	}
}

func (l *Logger) rotateIfNeeded() {
	currDate := time_helper.GetTimeNow().Format(timeFormat)

	l.mu.RLock()
	sameDate := currDate == l.date
	l.mu.RUnlock()

	if sameDate {
		return
	}

	newLogger := New()
	l.mu.Lock()
	l.Logger = newLogger.Logger
	l.date = newLogger.date
	l.mu.Unlock()
}

func (l *Logger) Info(entry global_dto.LogEntry) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.Logger.Info(entry.Message,
		zap.String("method", entry.Method),
		zap.String("endpoint", entry.Endpoint),
		zap.String("context", entry.Context),
		zap.String("scope", entry.Scope),
		zap.String("requestId", entry.RequestId),
		zap.String("startTime", entry.StartTime),
		zap.String("endTime", entry.EndTime),
		zap.Any("data", entry.Data),
	)
}

func (l *Logger) Warn(entry global_dto.LogEntry) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.Logger.Warn(entry.Message,
		zap.String("method", entry.Method),
		zap.String("endpoint", entry.Endpoint),
		zap.String("context", entry.Context),
		zap.String("scope", entry.Scope),
		zap.String("requestId", entry.RequestId),
		zap.String("startTime", entry.StartTime),
		zap.String("endTime", entry.EndTime),
		zap.String("stacktrace", entry.Stacktrace),
		zap.Any("data", entry.Data),
	)
}

func (l *Logger) Error(entry global_dto.LogEntry) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.Logger.Error(entry.Message,
		zap.String("method", entry.Method),
		zap.String("endpoint", entry.Endpoint),
		zap.String("context", entry.Context),
		zap.String("scope", entry.Scope),
		zap.String("requestId", entry.RequestId),
		zap.String("startTime", entry.StartTime),
		zap.String("endTime", entry.EndTime),
		zap.String("stacktrace", entry.Stacktrace),
		zap.Any("data", entry.Data),
	)
}

func (l *Logger) Debug(entry global_dto.LogEntry) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.Logger.Debug(entry.Message,
		zap.String("method", entry.Method),
		zap.String("endpoint", entry.Endpoint),
		zap.String("context", entry.Context),
		zap.String("scope", entry.Scope),
		zap.String("requestId", entry.RequestId),
		zap.String("startTime", entry.StartTime),
		zap.String("endTime", entry.EndTime),
		zap.Any("data", entry.Data),
	)
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
