package comps

import (
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

// UnavailableView is the unavailable component name.
const UnavailableView = "unavailable"

// Component represents a component of the patient view.
type Component interface {
	Init(moduleID string)

	OnConnect(client *ws.Client)
	OnDisconnect()

	OnModuleConnect()
	OnModuleDisconnect()

	Item() *Item
}

// Category represents a category displayed on the panel menu.
type Category struct {
	*js.Object
	Heading    string  `js:"heading"`
	SubHeading string  `js:"subHeading"`
	Icon       string  `js:"icon"`
	Items      []*Item `js:"items"`
}

// Item represents an item displayed in the menu as part of a category.
type Item struct {
	*js.Object
	Name       string `js:"name"`
	Heading    string `js:"heading"`
	SubHeading string `js:"subHeading"`
	Icon       string `js:"icon"`
	Component  string `js:"component"`
	Available  bool   `js:"available"`
	Active     bool   `js:"active"`
}
