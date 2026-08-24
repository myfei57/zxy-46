package plant

import (
	"errors"
	"time"

	"bldghvac/internal/audit"
	"bldghvac/internal/zone"
)

func (s *Service) RecoverBaseline(rec zone.SetbackRecord, now time.Time) error {
	if rec.Active {
		return errors.New("setback still active")
	}
	s.setpoint.SetBaseline(rec.DayBaseline)
	_ = s.audit.Record(audit.Entry{
		Source: "plant",
		ZoneID: rec.ZoneID,
		Action: "baseline-recover",
		Detail: "day baseline restored at " + now.Format(time.RFC3339),
	})
	return nil
}
