// +build !js

package lights

import (
	"strconv"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"
)

// HandleAction ...
func (m *Module) HandleAction(client *core.Client, act Action) error {
	arduino.Adaptor.PwmWrite(strconv.Itoa(m.Settings.Pin), byte(act.State))
	client.Emit(m.ID, Event{
		NewState: act.State,
	})

	return nil
}
