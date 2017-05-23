// +build js

package contact

import (
	"github.com/1lann/smarter-hospital/modules/ultrasonic"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/views/comps"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

const (
	contactDetectedIcon    = "green hotel"
	contactNotDetectedIcon = "grey help"
)

type Contact struct {
	moduleID  string
	item      *comps.Item
	component *ContactComponent
}

type ContactComponent struct {
	*js.Object
	Contact  bool `js:"contact"`
	moduleID string
}

func (c *Contact) Init(moduleID string) {
	item := &comps.Item{
		Object: js.Global.Get("Object").New(),
	}
	item.ID = moduleID
	item.Name = "Bed sensor"
	item.Component = "contact"
	item.Heading = "Not detected in bed"
	item.Icon = contactNotDetectedIcon
	item.Available = true
	item.Active = false

	component := &ContactComponent{
		Object:   js.Global.Get("Object").New(),
		moduleID: moduleID,
	}
	component.Contact = false

	views.ComponentWithTemplate(func() interface{} {
		return component
	}, "contact/contact.tmpl").Register("contact")

	c.moduleID = moduleID
	c.item = item
	c.component = component
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
	client.Subscribe(c.moduleID, c.onEvent)

	var info ultrasonic.Event
	err := views.ModuleInfo(c.moduleID, &info)
	if err != nil {
		println("contact info:", err.Error())
	}

	c.onEvent(info)
}

func (c *Contact) Item() *comps.Item {
	return c.item
}

func (c *Contact) OnDisconnect() {}

func (c *Contact) OnModuleConnect() {}

func (c *Contact) OnModuleDisconnect() {}
