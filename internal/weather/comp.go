package weather

import (
	"errors"
	"time"
)

type CompState string

const (
	CompNone    CompState = "none"
	CompPending CompState = "pending"
	CompApplied CompState = "applied"
)

type Compensation struct {
	State    CompState
	EffectAt time.Time
	Curve    float64
}

func (s *Service) PushCompensation(curve float64, effectAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.comp = Compensation{State: CompPending, EffectAt: effectAt, Curve: curve}
	return nil
}

func (s *Service) ApplyCompensation(now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.comp.State != CompPending {
		return nil
	}
	if now.Before(s.comp.EffectAt) {
		return errors.New("compensation not due")
	}
	s.comp.State = CompApplied
	return nil
}

func (s *Service) Compensation() Compensation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.comp
}
