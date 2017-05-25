// +build !js

package fsr

import (
	"strconv"
	"time"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"
)

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	pin := strconv.Itoa(m.Settings.Pin + 4)

	arduino.Adaptor.Board().On("AnalogRead"+pin, func(s interface{}) {
		r := m.getResistance((float64(s.(int)) / float64(1023)) * 5.0)

		pressed := false
		if r < m.Settings.Threshold {
			pressed = true
		}

		if pressed != m.LastEvent.Pressed {
			m.LastEvent = Event{
				Pressed: pressed,
			}
			client.Emit(m.ID, m.LastEvent)
		}
	})

	t := time.NewTicker(time.Second * 5)
	for range t.C {
		client.Emit(m.ID, m.LastEvent)
	}
}

func (m *Module) getResistance(v float64) float64 {
	i := (m.Settings.SupplyVoltage - v) / m.Settings.FixedResistor
	return (m.Settings.SupplyVoltage - i*m.Settings.FixedResistor) / i
}
