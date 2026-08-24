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
	unit, err := s.Chiller(id)
	if err != nil {
		return StartReport{}, err
	}
	if _, err := s.ElectLead(id); err != nil {
		return StartReport{}, err
	}
	unit.BootCompressor()
	pump, err := s.Pump(id)
	if err != nil {
		return StartReport{}, err
	}
	if err := chiller.StageLoad(unit, &pump); err != nil {
		return StartReport{}, err
	}
	s.pumps[id] = pump
	_ = s.persistPumps()
	report := StartReport{ChillerID: id, PumpStarted: pump.Running, CompressorStarted: unit.CompressorRunning}
	detail := "staged load, outdoor " + strconv.FormatFloat(s.weather.Latest().Temp, 'f', 1, 64) + "c"
	_ = s.audit.Record(audit.Entry{Source: "plant", ZoneID: id, Action: "start", Detail: detail})
	return report, nil
}

func (s *Service) StopChiller(id string) error {
	unit, err := s.Chiller(id)
	if err != nil {
		return err
	}
	unit.Stop()
	if pump, err := s.Pump(id); err == nil {
		pump.Stop()
		s.pumps[pump.ID] = pump
		_ = s.persistPumps()
	}
	_ = s.audit.Record(audit.Entry{Source: "plant", ZoneID: id, Action: "stop", Detail: "manual stop"})
	return s.persistChillers()
}
