// +build !js

package proximity

import (
	"time"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/pi/drivers"
	"github.com/1lann/vl53l0x"
)

// Module ...
type Module struct {
	ID string
	Settings

	LastEvent  Event
	Calibrated bool
	Threshold  int
	Count      int
}

func init() {
	core.RegisterModule(Module{})
}

// HandleEvent ...
func (m *Module) HandleEvent(evt Event) {
	m.LastEvent = evt
}

// HandleAction ...
func (m *Module) HandleAction(client *core.Client, act Action) error {
	// TODO: correction for when patient moves in/out of bed.
	m.Count += act.Correction
	if m.Count < 0 {
		m.Count = 0
	}

	client.Emit(m.ID, Event{
		Count: m.Count,
	})

	return nil
}

// Info ...
func (m *Module) Info() Event {
	return m.LastEvent
}

const (
	motionNone = iota
	motionIn
	motionOut
)

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	if m.Threshold == 0 {
		m.Threshold = 2000
	}

	d := vl53l0x.NewDriver(drivers.I2CBus)
	var calibrationCount int
	motionState := motionNone
	motionInTicks := 0
	motionOutTicks := 0

	lastDist := 0

	ticker := time.NewTicker(time.Millisecond * 50)
	for range ticker.C {
		dist, err := d.Measure()
		if err != nil {
			continue
		}

		if !m.Calibrated {
			calibrationCount++

			if dist < m.Threshold {
				m.Threshold = dist
			}

			if calibrationCount >= 60 {
				m.Threshold -= 10
				m.Calibrated = true
			}

			continue
		}

		if dist < m.Threshold && motionState == motionNone {
			motionState = motionIn
			motionInTicks = 0
		}

		if dist >= m.Threshold && motionState != motionNone {
			if motionState == motionIn {
				// Nothing happened
				motionState = motionNone
			} else if motionState == motionOut {
				motionOutTicks++

				if motionOutTicks <= motionInTicks {
					if m.Count == 0 {
						m.Count++
					} else {
						m.Count--
					}
				} else {
					m.Count++
				}

				client.Emit(m.ID, Event{
					Count: m.Count,
				})

				motionState = motionNone
			}
		}

		if motionState == motionIn {
			if dist > lastDist {
				if lastDist < (m.Threshold - m.Settings.PersonHeight) {
					motionState = motionOut
					motionOutTicks = 0
				}
			}
		}

		switch motionState {
		case motionIn:
			motionInTicks++
		case motionOut:
			motionOutTicks++
		}

		lastDist = dist
	}
}
