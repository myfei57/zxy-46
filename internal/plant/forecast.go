package plant

type LoadStep struct {
	ChillerID string
	Order     int
	Kw        float64
}

func (s *Service) LoadForecast() []LoadStep {
	s.mu.Lock()
	defer s.mu.Unlock()
	steps := make([]LoadStep, 0, len(s.order))
	for index, id := range s.order {
		unit := s.chillers[id]
		steps = append(steps, LoadStep{ChillerID: id, Order: index + 1, Kw: unit.LoadKw})
	}
	return steps
}

func (s *Service) PeakLoad() float64 {
	steps := s.LoadForecast()
	total := 0.0
	for _, step := range steps {
		total += step.Kw
	}
	return total
}
