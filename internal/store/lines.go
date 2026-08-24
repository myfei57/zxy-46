package store

import (
	"os"
	"strings"
)

func (s *Store) AppendLine(rel, line string) error {
	parts := strings.Split(rel, "/")
	if err := s.EnsureDir(parts[:len(parts)-1]...); err != nil {
		return err
	}
	file, err := os.OpenFile(s.Path(parts...), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(line + "\n")
	return err
}

func (s *Store) ReadLines(rel string) ([]string, error) {
	data, err := os.ReadFile(s.Path(strings.Split(rel, "/")...))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	content := strings.TrimRight(string(data), "\n")
	if content == "" {
		return nil, nil
	}
	return strings.Split(content, "\n"), nil
}
