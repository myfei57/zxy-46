package verifycase

import (
	"sync"
	"testing"

	"bldghvac/internal/audit"
	"bldghvac/internal/chiller"
	"bldghvac/internal/plant"
	"bldghvac/internal/quota"
	"bldghvac/internal/store"
	"bldghvac/internal/weather"
)

func TestHvacConcurrentChillerArbitration(t *testing.T) {
	dir := t.TempDir()
	st := store.New(dir)
	ws := weather.NewService(st)
	as := audit.NewService(st)
	qs := quota.NewService(st)
	ps := plant.NewService(st, ws, as, qs)
	if err := ps.RegisterChiller(chiller.NewUnit("c1", "主机1")); err != nil {
		t.Fatal(err)
	}
	if err := ps.RegisterChiller(chiller.NewUnit("c2", "主机2")); err != nil {
		t.Fatal(err)
	}
	if err := ps.RegisterPump(chiller.Pump{ID: "p1", ChillerID: "c1"}); err != nil {
		t.Fatal(err)
	}
	if err := ps.RegisterPump(chiller.Pump{ID: "p2", ChillerID: "c2"}); err != nil {
		t.Fatal(err)
	}
	const workers = 16
	loaded := make(chan string, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		id := "c1"
		if i%2 == 1 {
			id = "c2"
		}
		go func(chillerID string) {
			defer wg.Done()
			if report, err := ps.StartChiller(chillerID); err == nil && report.CompressorStarted {
				loaded <- report.ChillerID
			}
		}(id)
	}
	wg.Wait()
	close(loaded)
	leads := map[string]bool{}
	for id := range loaded {
		leads[id] = true
	}
	if len(leads) != 1 {
		t.Fatalf("exactly one chiller must load as lead, got %d", len(leads))
	}
}
