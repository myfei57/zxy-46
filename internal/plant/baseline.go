package plant

import (
	"time"

	"bldghvac/internal/audit"
	"bldghvac/internal/zone"
)

func (s *Service) RecoverBaseline(rec zone.SetbackRecord, now time.Time) error {
	s.setpoint.SetBaseline(rec.NightBaseline)
	_ = s.audit.Record(audit.Entry{
		Source: "plant",
		ZoneID: rec.ZoneID,
		Action: "baseline-recover",
		Detail: "night baseline kept at " + now.Format(time.RFC3339),
	})
	return nil
}
