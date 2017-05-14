// +build !js

package heartrate

import (
	"fmt"
	"log"
	"strconv"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"
)

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	// For some reason pin 4 = analog pin 0
	pin := strconv.Itoa(4 + m.Settings.Pin)

	for {
		val, err := arduino.Adaptor.AnalogRead(pin)
		if err != nil {
			log.Println("heartrate read:", err)
			continue
		}

		fmt.Println(val)
	}
}
