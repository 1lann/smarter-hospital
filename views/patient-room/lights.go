// +build js

package patientroom

import (
	"github.com/1lann/smarter-hospital/modules/lights"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

const (
	lightsOffIcon = "grey idea"
	lightsOnIcon  = "yellow idea"
)

type Lights struct {
	ModuleID  string
	item      *Item
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

func init() {
	item := &Item{
		Object: js.Global.Get("Object").New(),
	}
	item.Name = "Lights"
	item.Component = "lights"
	item.Heading = "Lights off"
	item.Icon = lightsOffIcon
	item.Available = true
	item.Active = false

	component := &LightsComponent{
		Object:   js.Global.Get("Object").New(),
		moduleID: "light1",
	}
	component.State = 0
	component.DimmedState = 20
	component.OnState = 70

	views.ComponentWithTemplate(func() interface{} {
		return component
	}, "patient-room/lights.tmpl").Register("lights")

	lightsComponent = &Lights{
		ModuleID:  component.moduleID,
		item:      item,
		component: component,
	}
}

func (c *LightsComponent) SetState(state int) {
	go func(state int) {
		_, err := views.ModuleDo(c.moduleID, lights.Action{State: state})
		if err != nil {
			println("lights: set state:", err.Error())
		}
	}(state)
}

func (l *Lights) onEvent(evt lights.Event) {
	if evt.NewState >= l.component.OnState {
		l.item.Heading = "Lights on"
		l.item.Icon = lightsOnIcon
	} else if evt.NewState >= l.component.DimmedState {
		l.item.Heading = "Lights dimmed"
		l.item.Icon = lightsOnIcon
	} else {
		l.item.Heading = "Lights off"
		l.item.Icon = lightsOffIcon
	}

	l.component.State = evt.NewState
}

func (l *Lights) OnConnect(client *ws.Client) {
	client.Subscribe(l.ModuleID, l.onEvent)

	var info lights.Event
	err := views.ModuleInfo(l.ModuleID, &info)
	if err != nil {
		println("lights info:", err.Error())
	}

	l.onEvent(info)
}

func (l *Lights) OnModuleConnect() {
	l.item.Available = true
	if l.item.Active {
		pageModel.ViewComponent = l.item.Component
	}
}

func (l *Lights) OnModuleDisconnect() {
	l.item.Available = false
	if l.item.Active {
		pageModel.ViewComponent = "unavailable"
	}
}

func (l *Lights) Item() *Item {
	return l.item
}
