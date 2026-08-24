package weather

import "math"

const ReturnDewPointLimit = 12.0

func DewPoint(temp, humidity float64) float64 {
	return temp - (100-humidity)/5
}

func BoundaryFor(r Reading) float64 {
	return math.Min(ReturnDewPointLimit, r.DewPoint-3)
}

func (s *Service) observe(r Reading) {
}

func (s *Service) Boundary() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.boundary
}

func (s *Service) RefreshBoundary() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observe(s.latest)
	return s.boundary
}
