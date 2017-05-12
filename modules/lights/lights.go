package lights

import "github.com/1lann/smarter-hospital/core"

// Module ...
type Module struct {
	ID string
	Settings

	CurrentState int
}

// Settings ...
type Settings struct {
	Pin int
}

// Action ...
type Action struct {
	State int // 0 = off, 255 = full
}

// Event ...
type Event struct {
	NewState int // 0 = off, 255 = full
}

func init() {
	core.RegisterModule(Module{})
}

// HandleEvent ...
func (m *Module) HandleEvent(evt Event) {
	m.CurrentState = evt.NewState
}

// Info ...
func (m *Module) Info() int {
	return m.CurrentState
}
