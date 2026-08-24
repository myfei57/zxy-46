package verifycase

import (
	"testing"
	"time"

	"bldghvac/internal/audit"
	"bldghvac/internal/plant"
	"bldghvac/internal/quota"
	"bldghvac/internal/store"
	"bldghvac/internal/weather"
	"bldghvac/internal/zone"
)

func TestHvacNightSetbackRecovery(t *testing.T) {
	dir := t.TempDir()
	st := store.New(dir)
	ws := weather.NewService(st)
	as := audit.NewService(st)
	qs := quota.NewService(st)
	ps := plant.NewService(st, ws, as, qs)
	zs := zone.NewService(st)
	enter := time.Date(2026, 8, 25, 23, 0, 0, 0, time.UTC)
	if err := zs.EnterSetback("z1", enter, 6, 12); err != nil {
		t.Fatal(err)
	}
	ps.Setpoint().SetBaseline(6)
	exit := time.Date(2026, 8, 26, 7, 0, 0, 0, time.UTC)
	record, err := zs.ExitSetback("z1", exit)
	if err != nil {
		t.Fatal(err)
	}
	if err := ps.RecoverBaseline(record, exit); err != nil {
		t.Fatal(err)
	}
	if got := ps.Baseline(); got != 12 {
		t.Fatalf("supply water baseline must recover to the day value, got %v", got)
	}
}
