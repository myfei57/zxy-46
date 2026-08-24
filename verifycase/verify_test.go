package verifycase

import (
	"sync"
	"testing"

	"bldghvac/internal/audit"
	"bldghvac/internal/plant"
	"bldghvac/internal/quota"
	"bldghvac/internal/store"
	"bldghvac/internal/vav"
	"bldghvac/internal/weather"
)

func TestHvacConcurrentVavBaseline(t *testing.T) {
	dir := t.TempDir()
	st := store.New(dir)
	ws := weather.NewService(st)
	as := audit.NewService(st)
	qs := quota.NewService(st)
	ps := plant.NewService(st, ws, as, qs)
	vs := vav.NewService(st, ps.Setpoint())
	const initial = 10.0
	ps.Setpoint().SetBaseline(initial)
	const workers = 200
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			vs.AdjustBaseline(1)
		}()
	}
	wg.Wait()
	if got := ps.Baseline(); got != initial+workers {
		t.Fatalf("shared baseline must accumulate every adjustment, got %v", got)
	}
}
