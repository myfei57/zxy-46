package quota

func (s *Service) AggregateDaily(day string) (DailyTotal, error) {
	total, err := s.sumDay(day)
	if err != nil {
		return DailyTotal{}, err
	}
	_, _ = s.st.FillTrendGaps(day)
	return total, nil
}

func (s *Service) ReportTotal(day string) (DailyTotal, error) {
	return s.sumDay(day)
}

func (s *Service) sumDay(day string) (DailyTotal, error) {
	points, err := s.st.ReadTrend(day)
	if err != nil {
		return DailyTotal{}, err
	}
	total := 0.0
	for _, point := range points {
		total += point.Kwh
	}
	return DailyTotal{Day: day, Kwh: total, Points: len(points)}, nil
}
