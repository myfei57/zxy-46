package chiller

import "errors"

type Pump struct {
	ID        string
	ChillerID string
	Running   bool
	Blocked   bool
}

func (p *Pump) Start() error {
	if p.Blocked {
		return errors.New("primary pump blocked")
	}
	p.Running = true
	return nil
}

func (p *Pump) Stop() {
	p.Running = false
}
