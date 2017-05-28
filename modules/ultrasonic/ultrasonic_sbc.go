// +build !js

package ultrasonic

import (
	"log"
	"time"

	"github.com/1lann/smarter-hospital/core"
	"github.com/kidoman/embd"
)

// Module ...
type Module struct {
	ID string
	Settings

	LastEvent Event
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

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	triggerPin, err := embd.NewDigitalPin(m.Settings.TriggerPin)
	if err != nil {
		panic(err)
	}

	triggerPin.SetDirection(embd.Out)

	echoPin, err := embd.NewDigitalPin(m.Settings.EchoPin)
	if err != nil {
		panic(err)
	}

	echoPin.SetDirection(embd.In)

	lastThree := make([]float64, 3)
	lastContact := false
	var lastEmit time.Time

	ticker := time.NewTicker(time.Millisecond * 100)
	for range ticker.C {
		triggerPin.Write(embd.High)
		go func() {
			time.Sleep(time.Millisecond * 10)
			triggerPin.Write(embd.Low)
		}()

		response := make(chan bool)

		var startTime time.Time
		var duration time.Duration

		err = echoPin.Watch(embd.EdgeBoth, func(arg2 embd.DigitalPin) {
			if startTime.IsZero() {
				startTime = time.Now()
				return
			}

			duration = time.Since(startTime)

			select {
			case response <- true:
			default:
			}
		})
		if err != nil {
			panic(err)
		}

		select {
		case <-response:
			break
		case <-time.After(time.Millisecond * 500):
			break
		}

		echoPin.StopWatching()

		if duration == 0 {
			continue
		}

		durationSeconds := duration.Seconds() * 1000
		if durationSeconds > m.Settings.ContactThreshold {
			durationSeconds = m.Settings.ContactThreshold * 1.5
		}

		log.Println(duration)

		lastThree = append(lastThree[1:], durationSeconds)
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
