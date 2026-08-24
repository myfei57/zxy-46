package demand

type Tariff struct {
	Name     string
	Baseline float64
}

func (s *Service) SwitchTariff(t Tariff) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tariff = t
	return s.st.WriteJSON("demand/tariff.json", tariffRecord{Tariff: t})
}

func (s *Service) Baseline() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tariff.Baseline
}

func (s *Service) Tariff() Tariff {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tariff
}
