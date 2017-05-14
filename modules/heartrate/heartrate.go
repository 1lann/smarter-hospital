package heartrate

import "github.com/1lann/smarter-hospital/core"

// Module ...
type Module struct {
	ID string
	Settings
}

// Settings ...
type Settings struct {
	PeakThreshold int
	Pin           int
}

// Action ...
type Action struct{}

// Event ...
type Event struct {
	BPM float64
}

func init() {
	core.RegisterModule(Module{})
}

// HandleEvent ...
func (m *Module) HandleEvent(evt Event) {}
