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

func TestHvacScheduleZoneOffsetBaseline(t *testing.T) {
	dir := t.TempDir()
	st := store.New(dir)
	zs := zone.NewService(st)
	ss := schedule.NewService(st)
	ws := weather.NewService(st)
	as := ahu.NewService(st, zs, ss, ws)
	if err := zs.UpsertProfile(zone.Profile{ZoneID: "east", Setpoint: 22, Deadband: 1, MinDamper: 30, LocalShift: 15}); err != nil {
		t.Fatal(err)
	}
	if err := zs.UpsertProfile(zone.Profile{ZoneID: "west", Setpoint: 22, Deadband: 1, MinDamper: 30}); err != nil {
		t.Fatal(err)
	}
	baseline := time.Date(2026, 8, 25, 8, 0, 0, 0, time.UTC)
	if err := ss.Add(schedule.Schedule{ID: "s-east", ZoneID: "east", Start: baseline}); err != nil {
		t.Fatal(err)
	}
	if err := ss.Add(schedule.Schedule{ID: "s-west", ZoneID: "west", Start: baseline}); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 25, 7, 50, 0, 0, time.UTC)
	east, err := as.EvaluatePrestart("east", baseline, now)
	if err != nil {
		t.Fatal(err)
	}
	west, err := as.EvaluatePrestart("west", baseline, now)
	if err != nil {
		t.Fatal(err)
	}
	if east.ShouldStart != west.ShouldStart {
		t.Fatalf("zero-offset zones must prestart together, east=%v west=%v", east.ShouldStart, west.ShouldStart)
	}
	if !east.ShouldStart {
		t.Fatal("zones must prestart inside the 15-minute window")
	}
}
