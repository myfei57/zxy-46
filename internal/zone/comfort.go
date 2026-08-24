package zone

import "math"

type Comfort struct {
	ZoneID string
	Score  float64
	Label  string
}

func (s *Service) ComfortScore(zoneID string, supply float64) (Comfort, error) {
	d, err := s.ComputeSetpoint(zoneID, supply)
	if err != nil {
		return Comfort{}, err
	}
	p, err := s.Profile(zoneID)
	if err != nil {
		return Comfort{}, err
	}
	delta := math.Abs(d.Supply - p.Setpoint)
	score := math.Max(0, 100-delta*10)
	label := "ok"
	if d.Supply < p.Setpoint-2 {
		label = "cold"
	} else if d.Supply > p.Setpoint+2 {
		label = "hot"
	}
	return Comfort{ZoneID: zoneID, Score: score, Label: label}, nil
}
