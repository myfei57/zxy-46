package vav

import (
	"errors"

	"bldghvac/internal/zone"
)

type Step struct {
	Order  int
	Action string
	Value  float64
}

type Sequence struct {
	ZoneID string
	Steps  []Step
}

func (s *Service) ExecuteSequence(zoneID string, d zone.Decision) (Sequence, error) {
	term, ok := s.terminalForZone(zoneID)
	if !ok {
		return Sequence{}, errors.New("no terminal for zone")
	}
	seq := Sequence{ZoneID: zoneID, Steps: []Step{}}
	if d.Reheat {
		seq.Steps = append(seq.Steps, Step{Order: 1, Action: "reheat-open", Value: 1})
		term.ReheatOpen = true
		seq.Steps = append(seq.Steps, Step{Order: 2, Action: "damper-set", Value: d.Damper})
		term.Damper = d.Damper
	} else {
		seq.Steps = append(seq.Steps, Step{Order: 1, Action: "damper-set", Value: d.Damper})
		term.Damper = d.Damper
		seq.Steps = append(seq.Steps, Step{Order: 2, Action: "reheat-close", Value: 0})
		term.ReheatOpen = false
	}
	term.Flow = d.Damper * 2
	_ = s.persist()
	return seq, nil
}
