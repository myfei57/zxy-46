package chiller

type State string

const (
	StateStandby  State = "standby"
	StatePumping  State = "pumping"
	StateLoading  State = "loading"
	StateRunning  State = "running"
	StateTripped  State = "tripped"
)

type Unit struct {
	ID                string
	Name              string
	State             State
	PumpRunning       bool
	CompressorRunning bool
	LoadKw            float64
	Tripped           bool
	TripReason        string
}

func NewUnit(id, name string) Unit {
	return Unit{ID: id, Name: name, State: StateStandby}
}

func (u *Unit) Stop() {
	u.State = StateStandby
	u.CompressorRunning = false
	u.LoadKw = 0
}

func (u *Unit) BootCompressor() {
	u.CompressorRunning = true
	u.LoadKw = 180
}
