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

outside:
	for {
		time.Sleep(time.Millisecond * 100)
		arduino.Adaptor.DigitalWrite(triggerPin, 1)
		time.Sleep(time.Millisecond)
		arduino.Adaptor.DigitalWrite(triggerPin, 0)

		riseStart := time.Now()

		for {
			result, err := arduino.Adaptor.DigitalRead(echoPin)
			if err != nil {
				panic(err)
			}

			if result == 1 {
				break
			}

			if time.Since(riseStart) > time.Second {
				log.Println("ultrasonic: rise timeout")
				continue outside
			}
		}

		pulseStart := time.Now()

		for {
			result, err := arduino.Adaptor.DigitalRead(echoPin)
			if err != nil {
				panic(err)
			}

			if result == 0 {
				break
			}

			if time.Since(pulseStart) > time.Second {
				log.Println("ultrasonic: fall timeout")
				continue outside
			}
		}

		lastThree = append(lastThree[1:], time.Now().Sub(pulseStart).Seconds()*1000)
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
