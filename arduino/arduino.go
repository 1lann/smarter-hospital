package arduino

import (
	"gobot.io/x/gobot/platforms/firmata"
)

// Adaptor represents the adaptor of the Firmata connection to the Arduino.
var Adaptor *firmata.Adaptor

// Connect starts a connection with the Arduino using the given serial
// connection.
func Connect(path string) error {
	Adaptor = firmata.NewAdaptor(path)
	return Adaptor.Connect()
}
