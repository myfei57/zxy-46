package plant

import (
	"errors"

	"bldghvac/internal/audit"
)

// ElectLead arbitrates the lead (host) chiller role for id. Only one chiller
// may hold the role at a time. When two start commands race, exactly one
// caller wins the role and the other is rejected as a standby so that only
// the lead is allowed to load — preventing the dual-host / dual-load
// overcurrent condition. The arbitration is serialized inside the
// LeadRegistry, so concurrent starters cannot both observe a win.
func (s *Service) ElectLead(id string) (string, error) {
	leader, ok := s.arbit.Claim(id)
	if !ok {
		return leader, errors.New("another chiller is lead")
	}
	_ = s.audit.Record(audit.Entry{
		Source: "plant",
		ZoneID: id,
		Action: "elect-lead",
		Detail: "arbitrated as lead chiller",
	})
	return leader, nil
}

func (s *Service) LeadChiller() string {
	return s.arbit.Leader()
}

// ClearTrip resets a chiller's protection trip. If the tripped chiller
// happened to hold the lead role it releases it so a healthy standby can
// take over rather than staying pinned to a unit that can no longer run.
func (s *Service) ClearTrip(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	unit, err := s.chillerLocked(id)
	if err != nil {
		return err
	}
	unit.ClearTrip()
	if s.arbit.Leader() == id {
		s.arbit.Release(id)
	}
	return s.persistChillers()
}
