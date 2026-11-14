package repo

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	// Baris tanggal dibuat dengan nomor 1, atau dinaikkan satu bila sudah ada. Baris itu terkunci
	// sampai transaksi pemanggil selesai, sehingga pesanan bersamaan tidak mendapat nomor sama.
	upsertCounterQuery = `INSERT INTO queue_counters (queue_date, last_number) VALUES (?, 1)
		ON DUPLICATE KEY UPDATE last_number = last_number + 1`
	getCounterQuery = `SELECT last_number FROM queue_counters WHERE queue_date = ? FOR UPDATE`
)

// NextQueueNumber mengambil nomor antrean berikutnya untuk tanggal now (keputusan 56, 95, 96).
// Pemanggil mengirim time_helper.GetTimeNow() (zona Asia/Jakarta) sehingga nomor reset pukul 00:00 WIB.
// Wajib dipanggil di dalam transaksi pembuatan pesanan (tahap 4.7): nomor dipakai saat pesanan dibuat,
// tidak dikembalikan walau pesanan batal, dan ikut batal bila pesanan gagal disimpan (rollback).
func NextQueueNumber(tx *gorm.DB, now time.Time) (int, error) {
	return nextForDate(tx, now)
}

func nextForDate(tx *gorm.DB, day time.Time) (int, error) {
	date := queueDate(day)
	if err := tx.Exec(upsertCounterQuery, date).Error; err != nil {
		return 0, err
	}
	var number int
	if err := tx.Raw(getCounterQuery, date).Scan(&number).Error; err != nil {
		return 0, err
	}
	return number, nil
}

// queueDate: tanggal kalender pada zona waktu milik t (WIB bila dari time_helper.GetTimeNow()).
func queueDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatQueueNumber: tiga digit tanpa awalan (001, 002, ...). Di atas 999 tampil apa adanya.
func FormatQueueNumber(n int) string {
	return fmt.Sprintf("%03d", n)
}
