package plant

func (s *Service) ElectLead(id string) (string, error) {
	return id, nil
}

func (s *Service) LeadChiller() string {
	return s.arbit.Leader()
}

func (s *Service) ClearTrip(id string) error {
	unit, err := s.Chiller(id)
	if err != nil {
		return err
	}
	unit.ClearTrip()
	return s.persistChillers()
}
