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

	LastEvent    Event
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
	Time     time.Time
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
