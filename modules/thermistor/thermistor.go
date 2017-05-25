package thermistor

import "github.com/1lann/smarter-hospital/core"

// Module ...
type Module struct {
	ID string
	Settings

	LastEvent Event
}

// Settings ...
type Settings struct {
	Pin           int
	R25           float64
	FixedResistor float64
	SupplyVoltage float64
}

// Action ...
type Action struct{}

// Event ...
type Event struct {
	Temperature float64 // In Kelvin
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
