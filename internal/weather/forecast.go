package weather

func (s *Service) Forecast() []Reading {
	latest := s.Latest()
	if latest.At.IsZero() {
		return nil
	}
	out := make([]Reading, 0, 4)
	for i := 1; i <= 4; i++ {
		out = append(out, Reading{
			At:       latest.At.AddDate(0, 0, i),
			Temp:     latest.Temp + float64(i),
			Humidity: latest.Humidity - float64(i*2),
		})
	}
	return out
}
