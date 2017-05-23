// +build js

package navbar

import (
	"time"

	"github.com/1lann/smarter-hospital/views"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jQuery = jquery.NewJQuery

// Model represents the mode for the navigation bar.
type Model struct {
	*js.Object
	Name      string `js:"name"`
	Date      string `js:"date"`
	Time      string `js:"time"`
	Connected bool   `js:"connected"`
}

func init() {
	views.ComponentWithTemplate(func() interface{} {
		m := &Model{Object: js.Global.Get("Object").New()}
		m.Name = views.GetUser().FirstName + " " + views.GetUser().LastName

		m.Date = time.Now().Format("Monday, _2 Jan 2006")
		m.Time = time.Now().Format("3:04:05 PM")

		go func() {
			for _ = range time.Tick(time.Second) {
				m.Date = time.Now().Format("Monday, _2 Jan 2006")
				m.Time = time.Now().Format("3:04:05 PM")
			}
		}()

		return m
	}, "nurse-navbar/nurse_navbar.tmpl", "connected").
		Register("nurse-navbar")
}
