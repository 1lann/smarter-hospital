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
	Name      string `js:"name"`
	Greeting  string `js:"greeting"`
	PingText  string `js:"pingText"`
	LightOn   bool   `js:"lightOn"`
	Connected bool   `js:"connected"`
	lastPing  time.Time
}

// Page represents the room page.
type Page struct {
	views.Page
	model *Model
}

// Title ...
func (p *Page) Title() string {
	return "Patient Room"
}
