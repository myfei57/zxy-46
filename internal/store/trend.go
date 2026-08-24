package store

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type TrendPoint struct {
	At     time.Time
	Kw     float64
	Kwh    float64
	Filled bool
}

func (s *Store) TrendRel(day string) string {
	return "trend/" + day + ".jsonl"
}

func (s *Store) AppendTrend(point TrendPoint) error {
	data, err := json.Marshal(point)
	if err != nil {
		return err
	}
	return s.AppendLine(s.TrendRel(point.At.Format("2006-01-02")), string(data))
}

func (s *Store) ReadTrend(day string) ([]TrendPoint, error) {
	lines, err := s.ReadLines(s.TrendRel(day))
	if err != nil {
		return nil, err
	}
	points := make([]TrendPoint, 0, len(lines))
	for _, line := range lines {
		var point TrendPoint
		if err := json.Unmarshal([]byte(line), &point); err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	return points, nil
}

func (s *Store) WriteTrend(day string, points []TrendPoint) error {
	var builder strings.Builder
	for _, point := range points {
		data, err := json.Marshal(point)
		if err != nil {
			return err
		}
		builder.Write(data)
		builder.WriteByte('\n')
	}
	rel := s.TrendRel(day)
	parts := strings.Split(rel, "/")
	if err := s.EnsureDir(parts[:len(parts)-1]...); err != nil {
		return err
	}
	return os.WriteFile(s.Path(parts...), []byte(builder.String()), 0o644)
}

func (s *Store) FillTrendGaps(day string) (int, error) {
	points, err := s.ReadTrend(day)
	if err != nil {
		return 0, err
	}
	if len(points) < 2 {
		return 0, nil
	}
	filled := 0
	for i := 0; i < len(points)-1; i++ {
		cur := points[i]
		next := points[i+1]
		step := int(next.At.Sub(cur.At).Minutes())
		if step > 1 {
			filled += step - 1
		}
	}
	return filled, nil
}
