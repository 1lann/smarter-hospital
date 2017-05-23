// +build js

package heartrate

import (
	"strconv"

	"github.com/1lann/smarter-hospital/modules/heartrate"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/views/comps"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

const (
	heartRateGood    = "green heartbeat"
	heartRateBad     = "red heartbeat"
	heartRateMissing = "grey help"
)

type HeartRate struct {
	moduleID  string
	item      *comps.Item
	component *HeartRateComponent
}

type HeartRateComponent struct {
	*js.Object
	BPM      int  `js:"bpm"`
	Contact  bool `js:"contact"`
	moduleID string
}

func (c *HeartRate) Init(moduleID string) {
	item := &comps.Item{
		Object: js.Global.Get("Object").New(),
	}
	item.ID = moduleID
	item.Name = "Heart rate sensor"
	item.Component = "heartrate"
	item.Heading = "No heart rate detected"
	item.Icon = heartRateMissing
	item.Available = true
	item.Active = false

	component := &HeartRateComponent{
		Object:   js.Global.Get("Object").New(),
		moduleID: moduleID,
	}
	component.Contact = false
	component.BPM = 0

	views.ComponentWithTemplate(func() interface{} {
		return component
	}, "heartrate/heartrate.tmpl").Register("heartrate")

	c.moduleID = moduleID
	c.item = item
	c.component = component
}

func (c *HeartRate) onEvent(evt heartrate.Event) {
	c.component.BPM = int(evt.BPM)
	c.component.Contact = evt.Contact

	if evt.Contact {
		c.item.Heading = strconv.Itoa(c.component.BPM) + " BPM"
		c.item.Icon = heartRateGood
	} else {
		c.item.Heading = "No heart rate detected"
		c.item.Icon = heartRateMissing
	}
}

func (c *HeartRate) OnConnect(client *ws.Client) {
	client.Subscribe(c.moduleID, c.onEvent)
}

func (c *HeartRate) Item() *comps.Item {
	return c.item
}

func (c *HeartRate) OnDisconnect() {}

func (c *HeartRate) OnModuleConnect() {}

func (c *HeartRate) OnModuleDisconnect() {}
