package demand

import (
	"sync"

	"bldghvac/internal/store"
	"bldghvac/internal/zone"
)

type tariffRecord struct {
	Tariff Tariff
}

type Service struct {
	st         *store.Store
	zoneSvc    *zone.Service
	mu         sync.Mutex
	tariff     Tariff
	cached     float64
	haveCached bool
}

func NewService(st *store.Store, zoneSvc *zone.Service) *Service {
	s := &Service{st: st, zoneSvc: zoneSvc}
	var record tariffRecord
	if err := st.ReadJSON("demand/tariff.json", &record); err == nil {
		s.tariff = record.Tariff
		s.cached = record.Tariff.Baseline
		s.haveCached = true
	}
	return s
}
