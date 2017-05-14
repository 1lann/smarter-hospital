// +build js

package patientroom

import (
	"time"

	"github.com/1lann/smarter-hospital/views"
	_ "github.com/1lann/smarter-hospital/views/notify"
	_ "github.com/1lann/smarter-hospital/views/patient-navbar"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

var pageModel *Model

type Component interface {
	OnConnect(client *ws.Client)
	OnDisconnect()
	Item() *Item
}

func (m *Model) SelectComponent(component string) {
	for _, category := range m.Categories {
		for _, item := range category.Items {
			if item.Component == component {
				item.Active = true
			} else {
				item.Active = false
			}
		}
	}

	m.ViewComponent = component
}

func getGreeting() string {
	hour := time.Now().Hour()
	if hour < 5 {
		return "evening"
	} else if hour < 12 {
		return "morning"
	} else if hour < 18 {
		return "afternoon"
	}

	return "evening"
}

func addItem(category *Category, name, subHeading, icon, component,
	moduleID string) {
	item := &Item{
		Object: js.Global.Get("Object").New(),
	}

	item.Name = name
	item.Heading = name
	item.SubHeading = subHeading
	item.Icon = icon
	item.Component = component
	item.Available = true
	item.Active = false

	category.Items = append(category.Items, item)
}

func populateCategories(m *Model) {
	agendaCategory := &Category{
		Object: js.Global.Get("Object").New(),
	}

	agendaCategory.Heading = ""
	agendaCategory.SubHeading = ""
	agendaCategory.Icon = ""
	agendaCategory.Items = make([]*Item, 0)

	addItem(agendaCategory, "Your agenda", "", "calendar", "agenda", "")
	m.Categories = append(m.Categories, agendaCategory)

	roomControls := &Category{
		Object: js.Global.Get("Object").New(),
	}

	roomControls.Heading = "Room controls"
	roomControls.SubHeading = ""
	roomControls.Icon = "settings"
	roomControls.Items = make([]*Item, 0)

	roomControls.Items = append(roomControls.Items, lightsComponent.Item())
	m.Categories = append(m.Categories, roomControls)

	healthCategory := &Category{
		Object: js.Global.Get("Object").New(),
	}

	healthCategory.Heading = "Your health"
	healthCategory.SubHeading = ""
	healthCategory.Icon = "green plus"
	healthCategory.Items = make([]*Item, 0)

	healthCategory.Items = append(healthCategory.Items, contactComponent.Item())
	m.Categories = append(m.Categories, healthCategory)
}

func (p *Page) OnLoad() {
	pageModel = &Model{
		Object: js.Global.Get("Object").New(),
	}

	views.ComponentWithTemplate(func() interface{} {
		return js.Global.Get("Object").New()
	}, "patient-room/unavailable.tmpl", "connected").Register("unavailable")

	pageModel.Categories = make([]*Category, 0)
	populateCategories(pageModel)

	p.model = pageModel

	pageModel.Name = views.GetUser().FirstName + " " + views.GetUser().LastName
	pageModel.Greeting = getGreeting()
	pageModel.PingText = ""
	pageModel.LightOn = false
	pageModel.Connected = false
	pageModel.ViewComponent = ""

	go func() {
		for _ = range time.Tick(time.Minute) {
			pageModel.Greeting = getGreeting()
		}
	}()

	views.ModelWithTemplate(pageModel, "patient-room/patient_room.tmpl")
}

func (p *Page) OnUnload(client *ws.Client) {
	if client != nil {
		client.Unsubscribe("pong")
		client.Unsubscribe("lights1")
	}
}

func (p *Page) OnConnect(client *ws.Client) {
	if p.connected {
		js.Global.Get("location").Call("reload")
	}
	p.connected = true
	p.model.Connected = true

	client.Subscribe("moduleConnect", func(moduleID string) {
		switch moduleID {
		case lightsComponent.ModuleID:
			lightsComponent.OnModuleConnect()
		case contactComponent.ModuleID:
			contactComponent.OnModuleConnect()
		}
	})

	client.Subscribe("moduleDisconnect", func(moduleID string) {
		println("disconnect:", moduleID)
		switch moduleID {
		case lightsComponent.ModuleID:
			lightsComponent.OnModuleDisconnect()
		case contactComponent.ModuleID:
			contactComponent.OnModuleDisconnect()
		}
	})

	connectedModules, err := views.ConnectedModules()
	if err != nil {
		println("connected modules error:", err.Error())
	}

	if isInList(lightsComponent.ModuleID, connectedModules) {
		lightsComponent.OnModuleConnect()
	} else {
		lightsComponent.OnModuleDisconnect()
	}

	if isInList(contactComponent.ModuleID, connectedModules) {
		contactComponent.OnModuleConnect()
	} else {
		contactComponent.OnModuleDisconnect()
	}

	lightsComponent.OnConnect(client)
	contactComponent.OnConnect(client)
}

func isInList(item string, list []string) bool {
	for _, listItem := range list {
		if item == listItem {
			return true
		}
	}
	return false
}

func (p *Page) OnDisconnect() {
	p.model.Connected = false
}
