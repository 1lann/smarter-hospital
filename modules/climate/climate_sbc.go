// +build !js

package climate

import (
	"time"

	"github.com/1lann/smarter-hospital/core"
	"github.com/kidoman/embd"
)

// Module ...
type Module struct {
	ID string
	Settings

	LastEvent  Event
	CoolingPin embd.PWMPin
	HeatingPin embd.PWMPin
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

// HandleAction ...
func (m *Module) HandleAction(client *core.Client, act Action) error {
	if act.Intensity > 1 {
		act.Intensity = 1
	}

	if act.Intensity < 0 {
		act.Intensity = 0.2
	}

	switch act.State {
	case StateOff:
		m.CoolingPin.SetAnalog(0)
		m.HeatingPin.SetAnalog(0)
	case StateCooling:
		m.HeatingPin.SetAnalog(0)
		m.HeatingPin.SetAnalog(byte(act.Intensity * float64(m.Settings.MaxCooling)))
	case StateHeating:
		m.CoolingPin.SetAnalog(0)
		m.HeatingPin.SetAnalog(byte(act.Intensity * float64(m.Settings.MaxHeating)))
	default:
		m.CoolingPin.SetAnalog(0)
		m.HeatingPin.SetAnalog(0)
		act.State = StateOff
	}

	m.LastEvent = Event{
		State:     act.State,
		Intensity: act.Intensity,
	}

	client.Emit(m.ID, m.LastEvent)

	return nil
}

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	var err error
	m.CoolingPin, err = embd.NewPWMPin(m.Settings.CoolingPin)
	if err != nil {
		panic(err)
	}

	m.HeatingPin, err = embd.NewPWMPin(m.Settings.HeatingPin)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second * 5)
		client.Emit(m.ID, m.LastEvent)
	}
}
