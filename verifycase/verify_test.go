package verifycase

import (
	"testing"
	"time"

	"bldghvac/internal/ahu"
	"bldghvac/internal/schedule"
	"bldghvac/internal/store"
	"bldghvac/internal/weather"
	"bldghvac/internal/zone"
)

func TestHvacFreeCoolingDewpoint(t *testing.T) {
	dir := t.TempDir()
	st := store.New(dir)
	zs := zone.NewService(st)
	ws := weather.NewService(st)
	ss := schedule.NewService(st)
	as := ahu.NewService(st, zs, ss, ws)
	now := time.Now()
	if err := ws.PushReading(weather.Reading{At: now, Temp: 24, Humidity: 90}); err != nil {
		t.Fatal(err)
	}
	latest := weather.Reading{At: now.Add(time.Minute), Temp: 15, Humidity: 60}
	if err := ws.PushReading(latest); err != nil {
		t.Fatal(err)
	}
	verdict := as.FreeCoolingVerdict("z1")
	if verdict.Allowed {
		t.Fatal("free cooling must be denied when the fresh dew point boundary is too tight")
	}
}
