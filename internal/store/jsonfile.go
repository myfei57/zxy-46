package store

import (
	"encoding/json"
	"os"
	"strings"
)

func (s *Store) WriteJSON(rel string, value any) error {
	parts := strings.Split(rel, "/")
	if err := s.EnsureDir(parts[:len(parts)-1]...); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	target := s.Path(parts...)
	tmp := target + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, target)
}

func (s *Store) ReadJSON(rel string, value any) error {
	data, err := os.ReadFile(s.Path(strings.Split(rel, "/")...))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}
