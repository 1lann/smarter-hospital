// +build js

package lights

import (
	"github.com/1lann/smarter-hospital/modules/lights"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/views/comps"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

const (
	lightsOffIcon = "grey idea"
	lightsOnIcon  = "yellow idea"
)

type Lights struct {
	moduleID  string
	item      *comps.Item
	component *LightsComponent
}

type LightsComponent struct {
	*js.Object
	State       int `js:"state"`
	OnState     int `js:"onState"`
	DimmedState int `js:"dimmedState"`
	moduleID    string
}

var lightsComponent *Lights

func (c *Lights) Init(moduleID string) {
	item := &comps.Item{
		Object: js.Global.Get("Object").New(),
	}
	item.ID = moduleID
	item.Name = "Lights"
	item.Component = "lights"
	item.Heading = "Lights off"
	item.Icon = lightsOffIcon
	item.Available = true
	item.Active = false

	component := &LightsComponent{
		Object:   js.Global.Get("Object").New(),
		moduleID: moduleID,
	}
	component.State = 0
	component.DimmedState = 20
	component.OnState = 70

	views.ComponentWithTemplate(func() interface{} {
		return component
	}, "lights/lights.tmpl").Register("lights")

	c.moduleID = moduleID
	c.item = item
	c.component = component
}

func (c *LightsComponent) SetState(state int) {
	go func(state int) {
		_, err := views.ModuleDo(c.moduleID, lights.Action{State: state})
		if err != nil {
			println("lights: set state:", err.Error())
		}
	}(state)
}

func (c *Lights) onEvent(evt lights.Event) {
	if evt.NewState >= c.component.OnState {
		c.item.Heading = "Lights on"
		c.item.Icon = lightsOnIcon
	} else if evt.NewState >= c.component.DimmedState {
		c.item.Heading = "Lights dimmed"
		c.item.Icon = lightsOnIcon
	} else {
		c.item.Heading = "Lights off"
		c.item.Icon = lightsOffIcon
	}

	c.component.State = evt.NewState
}

func (c *Lights) OnConnect(client *ws.Client) {
	client.Subscribe(c.moduleID, c.onEvent)

	var info lights.Event
	err := views.ModuleInfo(c.moduleID, &info)
	if err != nil {
		println("lights info:", err.Error())
	}

	c.onEvent(info)
}

func (c *Lights) Item() *comps.Item {
	return c.item
}

func (c *Lights) OnDisconnect() {}

func (c *Lights) OnModuleConnect() {}

func (c *Lights) OnModuleDisconnect() {}
