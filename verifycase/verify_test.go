package verifycase

import (
	"testing"

	"bldghvac/internal/audit"
	"bldghvac/internal/chiller"
	"bldghvac/internal/plant"
	"bldghvac/internal/quota"
	"bldghvac/internal/store"
	"bldghvac/internal/weather"
)

func TestHvacChillerStagedLoadingOrder(t *testing.T) {
	dir := t.TempDir()
	st := store.New(dir)
	ws := weather.NewService(st)
	as := audit.NewService(st)
	qs := quota.NewService(st)
	ps := plant.NewService(st, ws, as, qs)
	if err := ps.RegisterChiller(chiller.NewUnit("c1", "主机1")); err != nil {
		t.Fatal(err)
	}
	if err := ps.RegisterPump(chiller.Pump{ID: "p1", ChillerID: "c1", Blocked: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := ps.StartChiller("c1"); err == nil {
		t.Fatal("start must fail when the primary pump cannot run")
	}
	unit, err := ps.Chiller("c1")
	if err != nil {
		t.Fatal(err)
	}
	if unit.CompressorRunning {
		t.Fatal("compressor must not start before the primary pump")
	}
}
