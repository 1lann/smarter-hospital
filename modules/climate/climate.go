package climate

import "github.com/1lann/smarter-hospital/core"

type State int

const (
	StateOff = iota
	StateCooling
	StateHeating
)

// Module ...
type Module struct {
	ID string
	Settings

	LastEvent Event
}

// Settings ...
type Settings struct {
	CoolingPin int
	HeatingPin int

	MaxCooling byte
	MaxHeating byte
}

// Action ...
type Action struct {
	State     State
	Intensity float64
}

// Event ...
type Event struct {
	State     State
	Intensity float64
}

func init() {
	core.RegisterModule(Module{})
}

// HandleEvent ...
func (m *Module) HandleEvent(evt Event) {
	m.LastEvent = evt
}

// Info ...
func (m *Module) Info() Event {
	return m.LastEvent
}
