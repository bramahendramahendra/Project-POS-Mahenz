package scheduler

import (
	"context"
	"fmt"
	"time"

	"pos_api/config"
	cash_drawer_service "pos_api/domain/cash_drawer/service"
	"pos_api/helper"
	log_helper "pos_api/helper/log"
	"pos_api/model"
	"pos_api/repository"
)

const cashDrawerSchedulerName = "cash_drawer_auto_close"

func StartCashDrawerScheduler(ctx context.Context, svc cash_drawer_service.CashDrawerServiceInterface) {
	log_helper.LogInfo(log_helper.FromBackground("Scheduler", cashDrawerSchedulerName, "Cash drawer scheduler aktif"))

	go func() {
		for {
			now := time.Now().In(config.Location)
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, config.Location)
			timer := time.NewTimer(next.Sub(now))

			select {
			case <-timer.C:
				runCashDrawerAutoClose(svc)
			case <-ctx.Done():
				timer.Stop()
				log_helper.LogInfo(log_helper.FromBackground("Scheduler", cashDrawerSchedulerName, "Cash drawer scheduler berhenti"))
				return
			}
		}
	}()
}

func runCashDrawerAutoClose(svc cash_drawer_service.CashDrawerServiceInterface) {
	start := time.Now()

	log_helper.LogInfo(log_helper.FromBackground("Scheduler", cashDrawerSchedulerName, "Menjalankan auto close kas"))

	count, err := svc.AutoCloseYesterday()
	durationMs := time.Since(start).Milliseconds()

	status := "success"
	var message string

	if err != nil {
		status = "failed"
		message = err.Error()
		entry := log_helper.FromBackground("Scheduler", cashDrawerSchedulerName, "Auto close kas gagal")
		entry.Data = map[string]any{"error": err.Error()}
		log_helper.LogError(entry)
	} else if count == 0 {
		message = "Tidak ada kas yang perlu ditutup"
		log_helper.LogInfo(log_helper.FromBackground("Scheduler", cashDrawerSchedulerName, message))
	} else {
		message = fmt.Sprintf("%d kas berhasil ditutup otomatis", count)
		entry := log_helper.FromBackground("Scheduler", cashDrawerSchedulerName, message)
		entry.Data = map[string]any{"count": count}
		log_helper.LogInfo(entry)
	}

	saveLogScheduler(cashDrawerSchedulerName, status, message, durationMs)
}

func saveLogScheduler(schedulerName, status, message string, durationMs int64) {
	msg := message
	logData := &model.LogSchedulerModel{
		Id:            helper.GenerateUniqueId(),
		SchedulerName: schedulerName,
		Status:        status,
		Message:       &msg,
		DurationMs:    &durationMs,
		ExecutedAt:    time.Now(),
	}
	if err := repository.LogSchedulerRepo.InsertLogScheduler(logData); err != nil {
		entry := log_helper.FromBackground("Scheduler", schedulerName, "Gagal menyimpan log scheduler ke database")
		entry.Data = map[string]any{"error": err.Error()}
		log_helper.LogWarn(entry)
	}
}
