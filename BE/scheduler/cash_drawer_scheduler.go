package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"pos_api/config"
	cash_drawer_service "pos_api/domain/cash_drawer/service"
	"pos_api/helper"
	"pos_api/model"
	"pos_api/repository"
)

const cashDrawerSchedulerName = "cash_drawer_auto_close"

func StartCashDrawerScheduler(ctx context.Context, svc cash_drawer_service.CashDrawerServiceInterface) {
	log.Println("[Scheduler] Cash drawer scheduler aktif.")
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
				log.Println("[Scheduler] Cash drawer scheduler berhenti.")
				return
			}
		}
	}()
}

func runCashDrawerAutoClose(svc cash_drawer_service.CashDrawerServiceInterface) {
	start := time.Now()

	log.Println("[Scheduler] Menjalankan auto close kas...")

	count, err := svc.AutoCloseYesterday()
	durationMs := time.Since(start).Milliseconds()

	status := "success"
	var message string
	if err != nil {
		status = "failed"
		message = err.Error()
		log.Printf("[Scheduler] Auto close kas gagal: %v\n", err)
	} else if count == 0 {
		message = "Tidak ada kas yang perlu ditutup"
		log.Println("[Scheduler] Auto close kas selesai. Tidak ada kas yang perlu ditutup.")
	} else {
		message = fmt.Sprintf("%d kas berhasil ditutup otomatis", count)
		log.Printf("[Scheduler] Auto close kas selesai. %s.\n", message)
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
	_ = repository.LogSchedulerRepo.InsertLogScheduler(logData)
}
