package plant

import (
	"errors"

	"bldghvac/internal/audit"
	"bldghvac/internal/weather"
)

type SetpointState struct {
	BaselineValue float64
	compensated   float64
	compActive    bool
}

func NewSetpointState(baseline float64) *SetpointState {
	return &SetpointState{BaselineValue: baseline, compensated: baseline}
}

func (s *SetpointState) Baseline() float64 {
	return s.BaselineValue
}

func (s *SetpointState) SetBaseline(v float64) {
	s.BaselineValue = v
	s.compensated = v
	s.compActive = false
}

func (s *SetpointState) AdjustBaseline(delta float64) float64 {
	s.BaselineValue += delta
	s.compensated += delta
	return s.BaselineValue
}

func (s *SetpointState) ApplyCompensation(curve float64) {
	s.compActive = true
	s.compensated = curve
}

func (s *SetpointState) Current() float64 {
	if s.compActive {
		return s.compensated
	}
	return s.BaselineValue
}

func (s *Service) SyncWeatherComp(comp weather.Compensation) error {
	if comp.State != weather.CompApplied {
		return errors.New("compensation not applied")
	}
	s.setpoint.ApplyCompensation(comp.Curve)
	_ = s.audit.Record(audit.Entry{Source: "plant", ZoneID: "plant", Action: "weather-comp", Detail: "curve switched"})
	return nil
}
