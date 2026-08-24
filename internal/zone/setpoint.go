package zone

type Decision struct {
	Damper float64
	Reheat bool
	Supply float64
}

func (s *Service) ComputeSetpoint(zoneID string, supply float64) (Decision, error) {
	p, err := s.Profile(zoneID)
	if err != nil {
		return Decision{}, err
	}
	low := p.Setpoint - p.Deadband
	reheat := supply < low
	damper := p.MinDamper
	if !reheat && supply <= p.Setpoint+p.Deadband {
		damper = 100
	}
	return Decision{Damper: damper, Reheat: reheat, Supply: supply}, nil
}
