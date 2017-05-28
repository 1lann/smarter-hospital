// +build !js

package lights

import (
	"sync"
	"time"

	"github.com/1lann/smarter-hospital/core"
	"github.com/kidoman/embd"
)

// Module ...
type Module struct {
	ID string
	Settings

	LastEvent    Event
	CurrentState int
	Mutex        *sync.Mutex
	Pin          embd.PWMPin
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

// Init ...
func (m *Module) Init() {
	m.Mutex = new(sync.Mutex)
}

// HandleAction ...
func (m *Module) HandleAction(client *core.Client, act Action) error {
	go func(m *Module, act Action) {
		m.Mutex.Lock()
		defer m.Mutex.Unlock()

		if act.State == m.CurrentState {
			return
		}

		if m.CurrentState < act.State {
			ticks := act.State - m.CurrentState
			duration := m.AnimationDuration / time.Duration(ticks)
			for i := 0; i <= ticks; i++ {
				m.Pin.SetAnalog(byte(m.CurrentState + i))
				time.Sleep(duration)
			}
		} else {
			ticks := m.CurrentState - act.State
			duration := m.AnimationDuration / time.Duration(ticks)
			for i := 0; i <= ticks; i++ {
				m.Pin.SetAnalog(byte(m.CurrentState - i))
				time.Sleep(duration)
			}
		}

		m.Pin.SetAnalog(byte(act.State))

		m.CurrentState = act.State
	}(m, act)

	m.LastEvent = Event{
		NewState: act.State,
	}

	client.Emit(m.ID, m.LastEvent)

	return nil
}

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	var err error
	m.Pin, err = embd.NewPWMPin(m.Settings.Pin)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second * 5)
		client.Emit(m.ID, m.LastEvent)
	}
}
