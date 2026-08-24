package verifycase

import (
	"math"
	"testing"
	"time"

	"bldghvac/internal/quota"
	"bldghvac/internal/store"
)

func TestHvacTrendGapFillOrder(t *testing.T) {
	dir := t.TempDir()
	st := store.New(dir)
	base := time.Date(2026, 8, 25, 9, 0, 0, 0, time.UTC)
	if err := st.AppendTrend(store.TrendPoint{At: base, Kw: 100, Kwh: 2}); err != nil {
		t.Fatal(err)
	}
	if err := st.AppendTrend(store.TrendPoint{At: base.Add(5 * time.Minute), Kw: 120, Kwh: 4}); err != nil {
		t.Fatal(err)
	}
	qs := quota.NewService(st)
	total, err := qs.AggregateDaily("2026-08-25")
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(total.Kwh-18) > 1e-6 {
		t.Fatalf("daily total must include backfilled points, got %v", total.Kwh)
	}
}
