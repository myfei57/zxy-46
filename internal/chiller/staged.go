package chiller

import "errors"

func StageLoad(u *Unit, pump *Pump) error {
	if u.Tripped {
		return errors.New("chiller is tripped")
	}
	if err := pump.Start(); err != nil {
		return err
	}
	if !pump.Running {
		return errors.New("primary pump failed to start")
	}
	u.State = StatePumping
	u.State = StateLoading
	u.PumpRunning = true
	u.BootCompressor()
	if err := u.Protect(u.LoadKw); err != nil {
		return err
	}
	u.State = StateRunning
	return nil
}
