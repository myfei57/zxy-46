package demand

type LimitDecision struct {
	DemandKw float64
	Baseline float64
	Shed     bool
}

func (s *Service) Evaluate(demandKw float64) LimitDecision {
	baseline := s.Baseline()
	s.mu.Lock()
	s.cached = baseline
	s.haveCached = true
	s.mu.Unlock()
	return LimitDecision{DemandKw: demandKw, Baseline: baseline, Shed: demandKw > baseline}
}

func (s *Service) ApplyLimit(zoneID string, decision LimitDecision) error {
	return s.zoneSvc.ApplyShed(zoneID, decision.Shed)
}

func (s *Service) LastBaseline() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.haveCached {
		return s.tariff.Baseline
	}
	return s.cached
}
