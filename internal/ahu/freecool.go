package ahu

type FreeCoolingVerdict struct {
	Allowed bool
	Reason  string
}

func (s *Service) FreeCoolingVerdict(zoneID string) FreeCoolingVerdict {
	latest := s.weatherSvc.Latest()
	boundary := s.weatherSvc.Boundary()
	allowed := latest.DewPoint < boundary
	verdict := FreeCoolingVerdict{Allowed: allowed}
	if allowed {
		if s.zoneSvc.InSetback(zoneID) {
			allowed = false
			verdict = FreeCoolingVerdict{Allowed: false, Reason: "zone in setback"}
		} else {
			verdict.Reason = "dry enough for free cooling"
		}
	} else {
		verdict.Reason = "dew point above boundary"
	}
	if unit, ok := s.unitForZone(zoneID); ok {
		unit.FreeCooling = allowed
		_ = s.persist()
	}
	return verdict
}
