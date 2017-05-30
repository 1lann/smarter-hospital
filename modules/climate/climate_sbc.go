// +build !js

package climate

import (
	"strconv"
	"time"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/pi/drivers"
)

// Module ...
type Module struct {
	ID string
	Settings

	LastEvent  Event
	CoolingPin string
	HeatingPin string
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
		drivers.GoBot.PwmWrite(m.CoolingPin, 0)
		drivers.GoBot.PwmWrite(m.HeatingPin, 0)
	case StateCooling:
		drivers.GoBot.PwmWrite(m.HeatingPin, 0)
		drivers.GoBot.PwmWrite(m.HeatingPin, byte(act.Intensity*float64(m.Settings.MaxCooling)))
	case StateHeating:
		drivers.GoBot.PwmWrite(m.CoolingPin, 0)
		drivers.GoBot.PwmWrite(m.HeatingPin, byte(act.Intensity*float64(m.Settings.MaxHeating)))
	default:
		drivers.GoBot.PwmWrite(m.CoolingPin, 0)
		drivers.GoBot.PwmWrite(m.HeatingPin, 0)
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
	m.CoolingPin = strconv.Itoa(m.Settings.CoolingPin)
	m.HeatingPin = strconv.Itoa(m.Settings.HeatingPin)

	for {
		time.Sleep(time.Second * 5)
		client.Emit(m.ID, m.LastEvent)
	}
}
