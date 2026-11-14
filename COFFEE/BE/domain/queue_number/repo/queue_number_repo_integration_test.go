package repo

import (
	"sort"
	"sync"
	"testing"
	"time"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Uji nomor antrean (tahap 4.6) memakai DB dev Coffee nyata karena yang diuji adalah penguncian baris.
// Tanggal uji dibuat unik di tahun 1990-an agar tidak pernah bertabrakan dengan tanggal pesanan sungguhan.
// Baris uji dibiarkan di DB dev. Test dilewati bila MySQL dev tidak bisa dihubungi.
func connectCoffeeDevDBOrSkip(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "root:@tcp(127.0.0.1:3306)/pos_coffee_db?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Skipf("skip: DB dev tidak bisa dihubungi (%v)", err)
	}
	sqlDB, err := db.DB()
	if err != nil || sqlDB.Ping() != nil {
		t.Skip("skip: DB dev tidak bisa di-ping")
	}
	return db
}

// freshDay mencari tanggal 1990-an yang belum pernah dipakai (beserta hari berikutnya).
func freshDay(t *testing.T, db *gorm.DB) time.Time {
	t.Helper()
	day := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, int(time.Now().UnixNano()%3000))
	for i := 0; i < 5000; i++ {
		var n int64
		if err := db.Raw(`SELECT COUNT(*) FROM queue_counters WHERE queue_date IN (?, ?)`, day.Format("2006-01-02"), day.AddDate(0, 0, 1).Format("2006-01-02")).Scan(&n).Error; err != nil {
			t.Fatal(err)
		}
		if n == 0 {
			return day
		}
		day = day.AddDate(0, 0, 2)
	}
	t.Fatal("tidak menemukan tanggal uji yang kosong")
	return day
}

func TestNextQueueNumber_ConcurrentUniqueAndSequential(t *testing.T) {
	db := connectCoffeeDevDBOrSkip(t)
	day := freshDay(t, db)

	const workers = 20
	var wg sync.WaitGroup
	var mu sync.Mutex
	got := make([]int, 0, workers)
	errCh := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := db.Transaction(func(tx *gorm.DB) error {
				n, err := nextForDate(tx, day)
				if err != nil {
					return err
				}
				mu.Lock()
				got = append(got, n)
				mu.Unlock()
				return nil
			})
			if err != nil {
				errCh <- err
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatalf("gagal mengambil nomor: %v", err)
	}
	sort.Ints(got)
	for i, n := range got {
		if n != i+1 {
			t.Fatalf("nomor harus 1..%d tanpa ganda/lompat, didapat %v", workers, got)
		}
	}
}

func TestNextQueueNumber_NewDayStartsAtOne(t *testing.T) {
	db := connectCoffeeDevDBOrSkip(t)
	day := freshDay(t, db)
	next := day.AddDate(0, 0, 1)
	take := func(d time.Time) int {
		var n int
		if err := db.Transaction(func(tx *gorm.DB) error {
			var err error
			n, err = nextForDate(tx, d)
			return err
		}); err != nil {
			t.Fatalf("gagal mengambil nomor: %v", err)
		}
		return n
	}
	a, b := take(day), take(day)
	c := take(next)
	if a != 1 || b != 2 || c != 1 {
		t.Fatalf("hari pertama harus 1,2 lalu hari berikutnya 1; didapat %d,%d,%d", a, b, c)
	}
}

func TestNextQueueNumber_RolledBackTransactionDoesNotConsume(t *testing.T) {
	db := connectCoffeeDevDBOrSkip(t)
	day := freshDay(t, db)
	_ = db.Transaction(func(tx *gorm.DB) error {
		if _, err := nextForDate(tx, day); err != nil {
			return err
		}
		return gorm.ErrInvalidTransaction // paksa rollback (pesanan gagal disimpan)
	})
	var n int
	if err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		n, err = nextForDate(tx, day)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("pesanan yang gagal disimpan (rollback) tidak boleh memakai nomor; didapat %d", n)
	}
}

func TestFormatQueueNumber(t *testing.T) {
	cases := map[int]string{1: "001", 42: "042", 999: "999", 1000: "1000"}
	for in, want := range cases {
		if got := FormatQueueNumber(in); got != want {
			t.Errorf("FormatQueueNumber(%d) = %s, mau %s", in, got, want)
		}
	}
}

// Batas hari mengikuti WIB: 23:59 WIB masih tanggal itu, 00:00 WIB sudah tanggal berikutnya
// (padahal di UTC keduanya masih tanggal sebelumnya).
func TestQueueDate_ResetsAtMidnightWIB(t *testing.T) {
	wib := time.FixedZone("WIB", 7*3600)
	before := time.Date(2026, 11, 14, 16, 59, 0, 0, time.UTC).In(wib)
	after := time.Date(2026, 11, 14, 17, 0, 0, 0, time.UTC).In(wib)
	if got := queueDate(before); got != "2026-11-14" {
		t.Errorf("23:59 WIB harus 2026-11-14, didapat %s", got)
	}
	if got := queueDate(after); got != "2026-11-15" {
		t.Errorf("00:00 WIB harus 2026-11-15, didapat %s", got)
	}
}
