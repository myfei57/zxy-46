package plant

import (
	"errors"
	"sync"

	"bldghvac/internal/audit"
	"bldghvac/internal/weather"
)

type SetpointState struct {
	mu          sync.Mutex
	baseline    float64
	compensated float64
	compActive  bool
}

func NewSetpointState(baseline float64) *SetpointState {
	return &SetpointState{baseline: baseline, compensated: baseline}
}

func (s *SetpointState) Baseline() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.baseline
}

func (s *SetpointState) SetBaseline(v float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.baseline = v
	s.compensated = v
	s.compActive = false
}

func (s *SetpointState) AdjustBaseline(delta float64) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.baseline += delta
	s.compensated += delta
	return s.baseline
}

func (s *SetpointState) ApplyCompensation(curve float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.compActive = true
	s.compensated = curve
}

func (s *SetpointState) Current() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.compActive {
		return s.compensated
	}
	return s.baseline
}

func (s *Service) SyncWeatherComp(comp weather.Compensation) error {
	if comp.State != weather.CompApplied {
		return errors.New("compensation not applied")
	}
	s.setpoint.ApplyCompensation(comp.Curve)
	_ = s.audit.Record(audit.Entry{Source: "plant", ZoneID: "plant", Action: "weather-comp", Detail: "curve switched"})
	return nil
}
