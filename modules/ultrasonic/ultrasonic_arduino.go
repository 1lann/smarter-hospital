// +build !js

package ultrasonic

import (
	"log"
	"strconv"
	"time"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"
)

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	triggerPin := strconv.Itoa(m.Settings.TriggerPin)
	echoPin := strconv.Itoa(m.Settings.EchoPin)

	lastThree := make([]float64, 3)
	lastContact := false
	var lastEmit time.Time

	for {
		time.Sleep(time.Millisecond * 100)
		arduino.Adaptor.DigitalWrite(triggerPin, 1)
		go func() {
			time.Sleep(time.Millisecond * 10)
			arduino.Adaptor.DigitalWrite(triggerPin, 0)
		}()

		log.Println("pulse in")
		duration, err := arduino.Adaptor.PulseIn(echoPin, 1, time.Millisecond*100)
		if err != nil {
			log.Println("echo fail:", err)
			continue
		}

		log.Println("echo:", duration)

		lastThree = append(lastThree[1:], duration.Seconds()*1000)
		if lastThree[0] == 0 || lastThree[1] == 0 || lastThree[2] == 0 {
			continue
		}

		average := (lastThree[0] + lastThree[1] + lastThree[2]) / 3.0

		newContact := false
		if average < m.Settings.ContactThreshold {
			newContact = true
		}

		if lastContact != newContact || time.Since(lastEmit) > time.Second*3 {
			client.Emit(m.ID, Event{Contact: newContact})
			lastContact = newContact
			lastEmit = time.Now()
		}
	}
}
