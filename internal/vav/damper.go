package vav

import "errors"

func (s *Service) SetDamper(zoneID string, pct float64) error {
	term, ok := s.terminalForZone(zoneID)
	if !ok {
		return errors.New("no terminal for zone")
	}
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	term.Damper = pct
	term.Flow = pct * 2
	return s.persist()
}
