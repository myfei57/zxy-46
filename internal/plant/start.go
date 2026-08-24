package plant

import (
	"strconv"

	"bldghvac/internal/audit"
	"bldghvac/internal/chiller"
)

type StartReport struct {
	ChillerID         string
	PumpStarted       bool
	CompressorStarted bool
}

func (s *Service) StartChiller(id string) (StartReport, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	unit, err := s.chillerLocked(id)
	if err != nil {
		return StartReport{}, err
	}
	leader, err := s.ElectLead(id)
	if err != nil {
		// Lost the arbitration: another chiller is lead. id stays in
		// standby and must not load, otherwise two chillers load at once
		// and the condenser current hits the protection limit.
		_ = s.audit.Record(audit.Entry{
			Source: "plant",
			ZoneID: id,
			Action: "standby",
			Detail: "deferred to lead chiller " + leader,
		})
		return StartReport{}, err
	}
	pump, err := s.pumpLocked(id)
	if err != nil {
		return StartReport{}, err
	}
	if err := chiller.StageLoad(unit, &pump); err != nil {
		return StartReport{}, err
	}
	s.pumps[id] = pump
	if err := s.persistPumps(); err != nil {
		return StartReport{}, err
	}
	report := StartReport{ChillerID: id, PumpStarted: pump.Running, CompressorStarted: unit.CompressorRunning}
	detail := "staged load, outdoor " + strconv.FormatFloat(s.weather.Latest().Temp, 'f', 1, 64) + "c"
	_ = s.audit.Record(audit.Entry{Source: "plant", ZoneID: id, Action: "start", Detail: detail})
	return report, nil
}

func (s *Service) StopChiller(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	unit, err := s.chillerLocked(id)
	if err != nil {
		return err
	}
	unit.Stop()
	// A stopped chiller cannot serve as lead. Surrender the role so a
	// healthy standby can take over instead of being blocked.
	if s.arbit.Leader() == id {
		s.arbit.Release(id)
	}
	if pump, err := s.pumpLocked(id); err == nil {
		pump.Stop()
		s.pumps[pump.ID] = pump
		if err := s.persistPumps(); err != nil {
			return err
		}
	}
	_ = s.audit.Record(audit.Entry{Source: "plant", ZoneID: id, Action: "stop", Detail: "manual stop"})
	return s.persistChillers()
}
