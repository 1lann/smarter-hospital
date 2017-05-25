package fsr

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
	SupplyVoltage float64
	FixedResistor float64
	Threshold     float64
}

// Action ...
type Action struct{}

// Event ...
type Event struct {
	Pressed bool
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
