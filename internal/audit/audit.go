package audit

import (
	"encoding/json"

	"bldghvac/internal/store"
	"github.com/google/uuid"
)

type Service struct {
	st *store.Store
}

func NewService(st *store.Store) *Service {
	return &Service{st: st}
}

func (s *Service) Record(entry Entry) error {
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return s.st.AppendLine("audit/events.jsonl", string(data))
}

func (s *Service) List(limit int) ([]Entry, error) {
	lines, err := s.st.ReadLines("audit/events.jsonl")
	if err != nil {
		return nil, err
	}
	out := make([]Entry, 0, len(lines))
	for _, line := range lines {
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, err
		}
		out = append(out, entry)
	}
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out, nil
}

func (s *Service) Clear() error {
	return s.st.Delete("audit/events.jsonl")
}
