package datetime

import "time"

const TimeZoneAsiaBangkok = "Asia/Bangkok"

var (
	defaultTimeZone = TimeZoneAsiaBangkok
	defaultLocation = time.UTC
)

func init() {
	location, err := time.LoadLocation(TimeZoneAsiaBangkok)
	if err != nil {
		return
	}
	defaultLocation = location
}

func SetDefaultTimeZone(timeZone string) error {
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		return err
	}

	defaultTimeZone = timeZone
	defaultLocation = location
	return nil
}

func DefaultTimeZone() string {
	return defaultTimeZone
}

func GetCurrentLocation() *time.Location {
	return defaultLocation
}

func GetCurrentDateTimeNow() time.Time {
	return time.Now().In(defaultLocation)
}
