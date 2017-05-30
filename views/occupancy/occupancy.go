// +build js

package occupancy

import (
	"strconv"

	"github.com/1lann/smarter-hospital/modules/proximity"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/views/comps"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

const (
	hasOccupancyIcon = "green users"
	noOccupancyIcon  = "grey help"
)

type Occupancy struct {
	moduleID  string
	item      *comps.Item
	component *OccupancyComponent
}

type OccupancyComponent struct {
	*js.Object
	Occupancy int `js:"occupancy"`
	moduleID  string
}

func (c *Occupancy) Init(moduleID string) {
	item := &comps.Item{
		Object: js.Global.Get("Object").New(),
	}
	item.ID = moduleID
	item.Name = "Room occupancy"
	item.Component = "occupancy"
	item.Heading = "Room vacant"
	item.Icon = noOccupancyIcon
	item.Available = true
	item.Active = false

	component := &OccupancyComponent{
		Object:   js.Global.Get("Object").New(),
		moduleID: moduleID,
	}
	component.Occupancy = 0

	views.ComponentWithTemplate(func() interface{} {
		return component
	}, "occupancy/occupancy.tmpl").Register("occupancy")

	c.moduleID = moduleID
	c.item = item
	c.component = component
}

func (c *Occupancy) onEvent(evt proximity.Event) {
	c.component.Occupancy = evt.Count

	if evt.Count == 0 {
		c.item.Heading = "Room vacant"
		c.item.Icon = noOccupancyIcon
	} else if evt.Count == 1 {
		c.item.Heading = "1 person in the room"
		c.item.Icon = hasOccupancyIcon
	} else {
		c.item.Heading = strconv.Itoa(evt.Count) + " people in the room"
		c.item.Icon = hasOccupancyIcon
	}
}

func (c *Occupancy) OnConnect(client *ws.Client) {
	client.Subscribe(c.moduleID, c.onEvent)

	var info proximity.Event
	err := views.ModuleInfo(c.moduleID, &info)
	if err != nil {
		println("occupancy info:", err.Error())
	}

	c.onEvent(info)
}

func (c *Occupancy) Item() *comps.Item {
	return c.item
}

func (c *Occupancy) OnDisconnect() {}

func (c *Occupancy) OnModuleConnect() {}

func (c *Occupancy) OnModuleDisconnect() {}
