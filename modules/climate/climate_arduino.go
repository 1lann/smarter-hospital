// +build !js

package climate

import (
	"strconv"
	"time"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"
)

// HandleAction ...
func (m *Module) HandleAction(client *core.Client, act Action) error {
	coolingPin := strconv.Itoa(m.CoolingPin)
	heatingPin := strconv.Itoa(m.HeatingPin)

	if act.Intensity > 1 {
		act.Intensity = 1
	}

	if act.Intensity < 0 {
		act.Intensity = 0.2
	}

	switch act.State {
	case StateOff:
		arduino.Adaptor.PwmWrite(coolingPin, 0)
	case StateCooling:
		arduino.Adaptor.PwmWrite(heatingPin, 0)
		arduino.Adaptor.PwmWrite(coolingPin, byte(act.Intensity*float64(m.Settings.MaxCooling)))
	case StateHeating:
		arduino.Adaptor.PwmWrite(coolingPin, 0)
		arduino.Adaptor.PwmWrite(heatingPin, byte(act.Intensity*float64(m.Settings.MaxHeating)))
	default:
		arduino.Adaptor.PwmWrite(coolingPin, 0)
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
	for {
		time.Sleep(time.Second * 5)
		client.Emit(m.ID, m.LastEvent)
	}
}
