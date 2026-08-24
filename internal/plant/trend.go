package plant

import (
	"bldghvac/internal/quota"
	"bldghvac/internal/store"
)

func (s *Service) RecordTrend(point store.TrendPoint) error {
	return s.st.AppendTrend(point)
}

func (s *Service) CloseDay(day string) (quota.DailyTotal, error) {
	return s.quota.AggregateDaily(day)
}
