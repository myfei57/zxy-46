package quota

func (s *Service) OverBudget(day string, budgetKwh float64) (bool, error) {
	total, err := s.ReportTotal(day)
	if err != nil {
		return false, err
	}
	return total.Kwh > budgetKwh, nil
}
