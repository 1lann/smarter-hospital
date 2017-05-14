package lights

import (
	"sync"
	"time"

	"github.com/1lann/smarter-hospital/core"
)

// Module ...
type Module struct {
	ID string
	Settings

	CurrentState int
	Mutex        *sync.Mutex
}

// Settings ...
type Settings struct {
	Pin               int
	AnimationDuration time.Duration
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
func (m *Module) Info() Event {
	return Event{
		NewState: m.CurrentState,
	}
}
