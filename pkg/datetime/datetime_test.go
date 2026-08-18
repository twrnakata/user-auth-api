package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSetDefaultTimeZone_AsiaBangkok_NowUsesProjectTimezone(t *testing.T) {
	err := SetDefaultTimeZone(TimeZoneAsiaBangkok)
	require.NoError(t, err)
	require.Equal(t, TimeZoneAsiaBangkok, DefaultTimeZone())

	currentDateTime := GetCurrentDateTimeNow()
	require.Equal(t, TimeZoneAsiaBangkok, currentDateTime.Location().String())

	name, offset := currentDateTime.Zone()
	require.NotEmpty(t, name)
	require.Equal(t, 7*60*60, offset)
}

func TestSetDefaultTimeZone_InvalidTimeZone_ReturnsError(t *testing.T) {
	err := SetDefaultTimeZone("Not/AZone")
	require.Error(t, err)
}

func TestGetCurrentDateTimeNow_MatchesWallClockInProjectTimezone(t *testing.T) {
	err := SetDefaultTimeZone(TimeZoneAsiaBangkok)
	require.NoError(t, err)

	before := time.Now().In(GetCurrentLocation()).Add(-time.Second)
	currentDateTime := GetCurrentDateTimeNow()
	after := time.Now().In(GetCurrentLocation()).Add(time.Second)

	require.True(t, !currentDateTime.Before(before) && !currentDateTime.After(after))
}
