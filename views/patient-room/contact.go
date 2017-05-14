// +build js

package patientroom

import (
	"github.com/1lann/smarter-hospital/modules/ultrasonic"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

const (
	contactDetectedIcon    = "green hotel"
	contactNotDetectedIcon = "grey help"
)

type Contact struct {
	ModuleID  string
	item      *Item
	component *ContactComponent
}

type ContactComponent struct {
	*js.Object
	Contact  bool `js:"contact"`
	moduleID string
}

var contactComponent *Contact

func init() {
	item := &Item{
		Object: js.Global.Get("Object").New(),
	}
	item.Name = "Bed sensor"
	item.Component = "contact"
	item.Heading = "Not detected in bed"
	item.Icon = contactNotDetectedIcon
	item.Available = true
	item.Active = false

	component := &ContactComponent{
		Object:   js.Global.Get("Object").New(),
		moduleID: "ultrasonic1",
	}
	component.Contact = false

	views.ComponentWithTemplate(func() interface{} {
		return component
	}, "patient-room/contact.tmpl").Register("contact")

	contactComponent = &Contact{
		ModuleID:  component.moduleID,
		item:      item,
		component: component,
	}
}

func (c *Contact) onEvent(evt ultrasonic.Event) {
	c.component.Contact = evt.Contact

	if evt.Contact {
		c.item.Heading = "Detected in bed"
		c.item.Icon = contactDetectedIcon
	} else {
		c.item.Heading = "Not detected in bed"
		c.item.Icon = contactNotDetectedIcon
	}

}

func (c *Contact) OnConnect(client *ws.Client) {
	client.Subscribe(c.ModuleID, c.onEvent)

	var info ultrasonic.Event
	err := views.ModuleInfo(c.ModuleID, &info)
	if err != nil {
		println("contact info:", err.Error())
	}

	c.onEvent(info)
}

func (c *Contact) OnModuleConnect() {
	c.item.Available = true
	if c.item.Active {
		pageModel.ViewComponent = c.item.Component
	}
}

func (c *Contact) OnModuleDisconnect() {
	c.item.Available = false
	if c.item.Active {
		pageModel.ViewComponent = "unavailable"
	}
}

func (c *Contact) Item() *Item {
	return c.item
}
