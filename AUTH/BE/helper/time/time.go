package time_helper

import (
	"auth_api/config"
	"time"
)

func GetTimeNow() time.Time {
	return time.Now().In(config.Location)
}

func GetTimeWithFormat() string {
	return GetTimeNow().Format(config.FormatTime)
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
