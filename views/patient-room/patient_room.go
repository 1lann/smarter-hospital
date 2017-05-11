package patientroom

import "github.com/1lann/smarter-hospital/views"

func init() {
	views.RegisterPage("/patient/room", new(Page))
}

// Page represents the room page.
type Page struct {
	views.Page
}

// Title ...
func (p *Page) Title() string {
	return "Patient Room"
}
