// +build js

package patientroom

import (
	"strconv"

	"github.com/1lann/smarter-hospital/modules/heartrate"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

const (
	heartRateGood    = "green heartbeat"
	heartRateBad     = "red heartbeat"
	heartRateMissing = "grey help"
)

type HeartRate struct {
	ModuleID  string
	item      *Item
	component *HeartRateComponent
}

type HeartRateComponent struct {
	*js.Object
	BPM      int  `js:"bpm"`
	Contact  bool `js:"contact"`
	moduleID string
}

var heartrateComponent *HeartRate

func init() {
	item := &Item{
		Object: js.Global.Get("Object").New(),
	}
	item.Name = "Heart rate sensor"
	item.Component = "heartrate"
	item.Heading = "No heart rate detected"
	item.Icon = heartRateMissing
	item.Available = true
	item.Active = false

	component := &HeartRateComponent{
		Object:   js.Global.Get("Object").New(),
		moduleID: "heartrate1",
	}
	component.Contact = false
	component.BPM = 0

	views.ComponentWithTemplate(func() interface{} {
		return component
	}, "patient-room/heartrate.tmpl").Register("heartrate")

	heartrateComponent = &HeartRate{
		ModuleID:  component.moduleID,
		item:      item,
		component: component,
	}
}

func (h *HeartRate) onEvent(evt heartrate.Event) {
	h.component.BPM = int(evt.BPM)
	h.component.Contact = evt.Contact

	if evt.Contact {
		h.item.Heading = strconv.Itoa(h.component.BPM) + " BPM"
		h.item.Icon = heartRateGood
	} else {
		h.item.Heading = "No heart rate detected"
		h.item.Icon = heartRateMissing
	}
}

func (h *HeartRate) OnConnect(client *ws.Client) {
	client.Subscribe(h.ModuleID, h.onEvent)
}

func (h *HeartRate) OnModuleConnect() {
	h.item.Available = true
	if h.item.Active {
		pageModel.ViewComponent = h.item.Component
	}
}

func (h *HeartRate) OnModuleDisconnect() {
	h.item.Available = false
	if h.item.Active {
		pageModel.ViewComponent = "unavailable"
	}
}

func (h *HeartRate) Item() *Item {
	return h.item
}
