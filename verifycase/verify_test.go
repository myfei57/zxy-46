package verifycase

import (
	"testing"
	"time"

	"bldghvac/internal/audit"
	"bldghvac/internal/plant"
	"bldghvac/internal/quota"
	"bldghvac/internal/store"
	"bldghvac/internal/weather"
)

func TestHvacWeatherCompAppliesOnEffect(t *testing.T) {
	dir := t.TempDir()
	st := store.New(dir)
	ws := weather.NewService(st)
	as := audit.NewService(st)
	qs := quota.NewService(st)
	ps := plant.NewService(st, ws, as, qs)
	effect := time.Now().Add(5 * time.Minute)
	if err := ws.PushCompensation(6, effect); err != nil {
		t.Fatal(err)
	}
	if err := ws.ApplyCompensation(effect.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	comp := ws.Compensation()
	if comp.State != weather.CompApplied {
		t.Fatalf("compensation must be applied, got %s", comp.State)
	}
	if err := ps.SyncWeatherComp(comp); err != nil {
		t.Fatal(err)
	}
	if got := ps.Setpoint().Current(); got != 6 {
		t.Fatalf("setpoint must follow the new curve after effect, got %v", got)
	}
}
