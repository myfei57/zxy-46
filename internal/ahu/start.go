package ahu

import (
	"errors"
	"time"

	"bldghvac/internal/schedule"
)

func (s *Service) EvaluatePrestart(zoneID string, baseline time.Time, now time.Time) (schedule.PrestartDecision, error) {
	sch, err := s.schSvc.ForZone(zoneID)
	if err != nil {
		return schedule.PrestartDecision{}, err
	}
	profile, err := s.zoneSvc.Profile(zoneID)
	if err != nil {
		return schedule.PrestartDecision{}, err
	}
	return schedule.ComputePrestart(profile, sch.OffsetMinutes, baseline, now), nil
}

func (s *Service) StartUnit(id string) error {
	unit, err := s.Unit(id)
	if err != nil {
		return err
	}
	unit.Running = true
	unit.Prestart = false
	return s.persist()
}

func (s *Service) StopUnit(id string) error {
	unit, err := s.Unit(id)
	if err != nil {
		return err
	}
	unit.Running = false
	unit.FreeCooling = false
	return s.persist()
}

func (s *Service) MarkPrestart(zoneID string, flag bool) error {
	unit, ok := s.unitForZone(zoneID)
	if !ok {
		return errors.New("no ahu unit for zone")
	}
	unit.Prestart = flag
	return s.persist()
}
