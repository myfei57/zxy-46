package weather

import (
	"sync"
	"time"

	"bldghvac/internal/store"
)

type Reading struct {
	At       time.Time
	Temp     float64
	Humidity float64
	DewPoint float64
}

type Service struct {
	st      *store.Store
	mu      sync.RWMutex
	latest  Reading
	history []Reading
	comp    Compensation
	boundary float64
}

func NewService(st *store.Store) *Service {
	return &Service{st: st, boundary: ReturnDewPointLimit}
}

func (s *Service) PushReading(r Reading) error {
	if r.DewPoint == 0 {
		r.DewPoint = DewPoint(r.Temp, r.Humidity)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latest = r
	s.history = append(s.history, r)
	if len(s.history) > 200 {
		s.history = s.history[len(s.history)-200:]
	}
	s.observe(r)
	return s.st.AppendLine("weather/readings.jsonl", "record")
}

func (s *Service) Latest() Reading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.latest
}

func (s *Service) History(n int) []Reading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if n <= 0 || n > len(s.history) {
		n = len(s.history)
	}
	out := make([]Reading, n)
	copy(out, s.history[len(s.history)-n:])
	return out
}
