package patientroom

import (
	"time"

	"github.com/1lann/smarter-hospital/views"
	"github.com/gopherjs/gopherjs/js"
)

func init() {
	views.RegisterPage("/patient/room", new(Page))
}

// Model represents the Vue.js model.
type Model struct {
	*js.Object
	Name          string      `js:"name"`
	Greeting      string      `js:"greeting"`
	PingText      string      `js:"pingText"`
	LightOn       bool        `js:"lightOn"`
	Connected     bool        `js:"connected"`
	Categories    []*Category `js:"categories"`
	ViewComponent string      `js:"viewComponent"`
	lastPing      time.Time
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

// Page represents the room page.
type Page struct {
	connected bool
	views.Page
	model *Model
}

// Title ...
func (p *Page) Title() string {
	return "Patient Room"
}
