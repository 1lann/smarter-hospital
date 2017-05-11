// +build js

package patientroom

import (
	"time"

	"github.com/1lann/smarter-hospital/views"
	_ "github.com/1lann/smarter-hospital/views/notify"
	_ "github.com/1lann/smarter-hospital/views/patient-navbar"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"

	"github.com/1lann/smarter-hospital/modules/ping"
)

type Model struct {
	*js.Object
	Name     string `js:"name"`
	Greeting string `js:"greeting"`
	PingText string `js:"pingText"`
	lastPing time.Time
}

func (m *Model) Ping() {
	go func() {
		m.lastPing = time.Now()
		text := m.PingText
		m.PingText = ""
		_, err := views.Do("ping1", ping.Action{
			Message: text,
		})
		if err != nil {
			println("not nil error:", err.Error())
			return
		}
		println("Successful ping!")
	}()
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

func (p *Page) OnLoad() {
	m := &Model{
		Object: js.Global.Get("Object").New(),
	}

	m.Name = views.GetUser().FirstName + " " + views.GetUser().LastName
	m.Greeting = getGreeting()
	m.PingText = ""

	go func() {
		for _ = range time.Tick(time.Minute) {
			m.Greeting = getGreeting()
		}
	}()

	views.ModelWithTemplate(m, "patient-room/patient_room.tmpl")

	client := ws.NewClient()
	client.HandleEvent("pong", func(msg string) {
		println("Received pong:", msg)
		println("Latency:", time.Since(m.lastPing).String())
	})
	// TODO: For the shits and giggles
	client.Connect("ws://127.0.0.1:8080/ws")
}
