// +build !js

package lights

import (
	"strconv"
	"sync"
	"time"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"
)

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
				arduino.Adaptor.PwmWrite(strconv.Itoa(m.Settings.Pin), byte(m.CurrentState+i))
				time.Sleep(duration)
			}
		} else {
			ticks := m.CurrentState - act.State
			duration := m.AnimationDuration / time.Duration(ticks)
			for i := 0; i <= ticks; i++ {
				arduino.Adaptor.PwmWrite(strconv.Itoa(m.Settings.Pin), byte(m.CurrentState-i))
				time.Sleep(duration)
			}
		}

		arduino.Adaptor.PwmWrite(strconv.Itoa(m.Settings.Pin), byte(act.State))

		m.CurrentState = act.State
	}(m, act)

	client.Emit(m.ID, Event{
		NewState: act.State,
		Time:     time.Now(),
	})

	return nil
}
