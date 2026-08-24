package audit

func (s *Service) BySource(source string) ([]Entry, error) {
	entries, err := s.List(0)
	if err != nil {
		return nil, err
	}
	out := make([]Entry, 0)
	for _, entry := range entries {
		if entry.Source == source {
			out = append(out, entry)
		}
	}
	return out, nil
}
