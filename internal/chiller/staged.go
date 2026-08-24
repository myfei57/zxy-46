package chiller

import "errors"

func StageLoad(u *Unit, pump *Pump) error {
	if u.Tripped {
		return errors.New("chiller is tripped")
	}
	u.State = StateLoading
	// 规程：先起导泵，再合主机。导泵未建流前不得合主机，否则负荷顶到冷凝器保护。
	if err := pump.Start(); err != nil {
		return err
	}
	u.PumpRunning = true
	u.BootCompressor()
	u.State = StateRunning
	return nil
}
