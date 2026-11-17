package time_helper

import (
	"pos_api/config"
	"time"
)

func GetTimeNow() time.Time {
	return time.Now().In(config.Location)
}

func GetTimeWithFormat() string {
	return GetTimeNow().Format(config.FormatTime)
}

// StartOfDay mengembalikan jam 00:00:00 WIB dari tanggal t.
func StartOfDay(t time.Time) time.Time {
	t = t.In(config.Location)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, config.Location)
}

// ToSQLDate memformat t menjadi "2006-01-02" (WIB) untuk dipakai sebagai parameter query MySQL.
func ToSQLDate(t time.Time) string {
	return t.In(config.Location).Format("2006-01-02")
}

func GetEndTime(timeString string) string {
	now, err := time.ParseInLocation(config.FormatTime, timeString, config.Location)
	if err != nil {
		panic("failed to parsing time string format, " + err.Error())
	}

	duration := time.Since(now)
	endTime := now.Add(duration)
	return endTime.Format(config.FormatTime)
}

// NormalizeDateRange melengkapi date_from/date_to yang cuma berisi tanggal (YYYY-MM-DD)
// dengan jam awal/akhir hari, supaya BETWEEN di query SQL mencakup seluruh hari tersebut.
// Kalau kosong, keduanya default ke hari ini.
func NormalizeDateRange(dateFrom, dateTo string) (string, string) {
	today := GetTimeNow().Format("2006-01-02")
	if dateFrom == "" {
		dateFrom = today
	}
	if dateTo == "" {
		dateTo = today
	}
	if len(dateFrom) == len("2006-01-02") {
		dateFrom += " 00:00:00"
	}
	if len(dateTo) == len("2006-01-02") {
		dateTo += " 23:59:59"
	}
	return dateFrom, dateTo
}

