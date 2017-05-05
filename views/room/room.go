package room

import "github.com/1lann/smarter-hospital/views"

func init() {
	views.RegisterPage("/room", new(Page))
}

// Page represents the room page.
type Page struct {
	views.Page
}

// Title ...
func (p *Page) Title() string {
	return "Room"
}
