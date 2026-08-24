package quota

import "strconv"

func (s *Service) RecordMeter(day string, kwh float64) error {
	return s.st.AppendLine("quota/meter-"+day+".jsonl", strconv.FormatFloat(kwh, 'f', 2, 64))
}

func (s *Service) MeterReadings(day string) ([]float64, error) {
	lines, err := s.st.ReadLines("quota/meter-" + day + ".jsonl")
	if err != nil {
		return nil, err
	}
	out := make([]float64, 0, len(lines))
	for _, line := range lines {
		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, err
		}
		out = append(out, value)
	}
	return out, nil
}
