package plant

import "errors"

func (s *Service) ElectLead(id string) (string, error) {
	lead, isLead := s.arbit.Claim(id)
	if !isLead {
		return lead, errors.New("not the lead chiller")
	}
	return lead, nil
}

func (s *Service) LeadChiller() string {
	return s.arbit.Leader()
}

func (s *Service) ClearTrip(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	unit, err := s.chillerLocked(id)
	if err != nil {
		return err
	}
	unit.ClearTrip()
	return s.persistChillers()
}
