package quota

import "fmt"

type DailyInsight struct {
	Day           string
	TotalKwh      float64
	MeterCount    int
	AverageKwh    float64
	FilledPoints  int
	Status        string
}

func (s *Service) DailyInsight(day string) (DailyInsight, error) {
	total, err := s.ReportTotal(day)
	if err != nil {
		return DailyInsight{}, err
	}
	meters, err := s.MeterReadings(day)
	if err != nil {
		return DailyInsight{}, err
	}
	average := 0.0
	if len(meters) > 0 {
		sum := 0.0
		for _, value := range meters {
			sum += value
		}
		average = sum / float64(len(meters))
	}
	status := "ok"
	if total.Points > 0 && total.Kwh > 0 {
		status = fmt.Sprintf("points:%d", total.Points)
	}
	return DailyInsight{
		Day:          day,
		TotalKwh:     total.Kwh,
		MeterCount:   len(meters),
		AverageKwh:   average,
		FilledPoints: total.Points,
		Status:       status,
	}, nil
}
