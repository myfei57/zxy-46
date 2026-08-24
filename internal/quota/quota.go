package quota

import "bldghvac/internal/store"

type DailyTotal struct {
	Day    string
	Kwh    float64
	Points int
}

type Service struct {
	st *store.Store
}

func NewService(st *store.Store) *Service {
	return &Service{st: st}
}
