// +build js

package room

import (
	"time"

	"github.com/1lann/smarter-hospital/views"
	_ "github.com/1lann/smarter-hospital/views/notify"
	_ "github.com/1lann/smarter-hospital/views/patient-navbar"
	"github.com/gopherjs/gopherjs/js"
)

type Model struct {
	*js.Object
	Message string `js:"message"`
}

func (p *Page) OnLoad() {
	println("Hello, world!")

	m := &Model{
		Object: js.Global.Get("Object").New(),
	}

	m.Message = "Hello web apps from Go!"

	views.ModelWithTemplate(m, "room/room.tmpl")

	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			m.Message = time.Now().String()
		}
	}()
}
