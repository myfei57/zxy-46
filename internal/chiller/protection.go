package chiller

import "errors"

const condenserLimitKw = 420

func (u *Unit) Protect(kw float64) error {
	if kw > condenserLimitKw {
		u.Tripped = true
		u.State = StateTripped
		u.TripReason = "condenser overcurrent"
		return errors.New("condenser protection tripped")
	}
	u.LoadKw = kw
	return nil
}

func (u *Unit) ClearTrip() {
	u.Tripped = false
	u.TripReason = ""
	u.State = StateStandby
}
