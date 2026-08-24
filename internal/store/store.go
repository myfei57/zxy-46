package store

import (
	"os"
	"path/filepath"
	"sort"
)

type Store struct {
	root string
}

func New(root string) *Store {
	return &Store{root: root}
}

func (s *Store) Root() string {
	return s.root
}

func (s *Store) Path(parts ...string) string {
	return filepath.Join(append([]string{s.root}, parts...)...)
}

func (s *Store) EnsureDir(parts ...string) error {
	return os.MkdirAll(s.Path(parts...), 0o755)
}

func (s *Store) Exists(parts ...string) bool {
	_, err := os.Stat(s.Path(parts...))
	return err == nil
}

func (s *Store) Delete(parts ...string) error {
	err := os.Remove(s.Path(parts...))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *Store) List(parts ...string) ([]string, error) {
	entries, err := os.ReadDir(s.Path(parts...))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}
