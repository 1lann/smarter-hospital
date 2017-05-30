// +build !js

package logic

import (
	"math"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/modules/climate"
	"github.com/1lann/smarter-hospital/modules/heartrate"
	"github.com/1lann/smarter-hospital/modules/lights"
	"github.com/1lann/smarter-hospital/modules/proximity"
	"github.com/1lann/smarter-hospital/modules/thermometer"
	"github.com/1lann/smarter-hospital/modules/ultrasonic"
	"github.com/1lann/smarter-hospital/notify"
)

func (c *ClimateControl) Handle(server *core.Server, clim climate.Module, therm thermometer.Module) {
	if therm.Info().Temperature == c.state.CurrentTemperature {
		return
	}

	c.state = ClimateState{
		On:                 c.state.On,
		State:              clim.Info().State,
		TargetTemperature:  c.state.TargetTemperature,
		CurrentTemperature: therm.Info().Temperature,
	}

	if !c.state.On {
		server.Do(clim.ID, climate.Action{
			State:     climate.StateOff,
			Intensity: 0,
		})

		return
	}

	if int(math.Ceil(therm.Info().Temperature-0.5)) < c.state.TargetTemperature {
		server.Do(clim.ID, climate.Action{
			State:     climate.StateHeating,
			Intensity: 1,
		})
	} else if int(math.Ceil(therm.Info().Temperature-0.5)) > c.state.TargetTemperature {
		server.Do(clim.ID, climate.Action{
			State:     climate.StateCooling,
			Intensity: 1,
		})
	} else {
		server.Do(clim.ID, climate.Action{
			State:     climate.StateOff,
			Intensity: 0,
		})
	}

	c.wsServer.Emit("climatecontrol", c.state)
}

func (w *Warner) Handle(server *core.Server, hr heartrate.Module, u ultrasonic.Module) {
	if !hr.Info().Contact {
		if !w.hasHeartrateAlert {
			w.hasHeartrateAlert = true
			w.notifyServer.Push(notify.Notification{
				Alert:      true,
				Dismissed:  false,
				Heading:    "Heart rate sensor removed",
				SubHeading: "Ash Ketchum - Room 025",
				Location:   "Ash Ketchum - Room 025",
				Icon:       "red heartbeat",
				Link:       "/nurse/room/#heartrate1",
			})
		}
	} else {
		if hr.Info().BPM < 50 {
			if !w.hasHeartrateAlert {
				w.hasHeartrateAlert = true
				w.notifyServer.Push(notify.Notification{
					Alert:      true,
					Dismissed:  false,
					Heading:    "Low heart rate",
					SubHeading: "Ash Ketchum - Room 025",
					Location:   "Ash Ketchum - Room 025",
					Icon:       "red heartbeat",
					Link:       "/nurse/room/#heartrate1",
				})
			}
		} else if hr.Info().BPM > 100 && !w.hasHeartrateAlert {
			if !w.hasHeartrateAlert {
				w.hasHeartrateAlert = true
				w.notifyServer.Push(notify.Notification{
					Alert:      true,
					Dismissed:  false,
					Heading:    "High heart rate",
					SubHeading: "Ash Ketchum - Room 025",
					Location:   "Ash Ketchum - Room 025",
					Icon:       "red heartbeat",
					Link:       "/nurse/room/#heartrate1",
				})
			}
		} else {
			w.hasHeartrateAlert = false
		}
	}

	if !u.Info().Contact {
		if !w.hasBedAlert {
			w.hasBedAlert = true
			if !w.hasBedAlert {
				w.notifyServer.Push(notify.Notification{
					Alert:      true,
					Dismissed:  false,
					Heading:    "Patient left bed",
					SubHeading: "Ash Ketchum - Room 025",
					Location:   "Ash Ketchum - Room 025",
					Icon:       "red warning sign",
					Link:       "/nurse/room/#ultrasonic1",
				})
			}
		}
	} else {
		w.hasBedAlert = false
	}
}

func (s *SmartLighting) Handle(server *core.Server, prox proximity.Module) {
	if prox.Info().Count == 0 {
		if !s.hasChanged {
			s.hasChanged = true
			server.Do("light1", lights.Action{
				State: 0,
			})
		}
	} else {
		s.hasChanged = false
	}
}
