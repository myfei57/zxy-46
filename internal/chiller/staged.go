package chiller

import "errors"

func StageLoad(u *Unit, pump *Pump) error {
	if u.Tripped {
		return errors.New("chiller is tripped")
	}
	u.State = StateLoading
	u.PumpRunning = true
	u.BootCompressor()
	u.State = StateRunning
	if err := pump.Start(); err != nil {
		return err
	}
	return nil
}
