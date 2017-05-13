// +build js

package patientroom

import (
	"time"

	"github.com/1lann/smarter-hospital/views"
	_ "github.com/1lann/smarter-hospital/views/notify"
	_ "github.com/1lann/smarter-hospital/views/patient-navbar"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"

	"github.com/1lann/smarter-hospital/modules/lights"
	"github.com/1lann/smarter-hospital/modules/ping"
)

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

func (m *Model) SetLight(num int) {
	go func() {
		println("Setting light to", num)
		_, err := views.Do("lights1", lights.Action{
			State: num,
		})
		if err != nil {
			println(":(", err.Error())
		}
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
	p.model = m

	m.Name = views.GetUser().FirstName + " " + views.GetUser().LastName
	m.Greeting = getGreeting()
	m.PingText = ""
	m.LightOn = false
	m.Connected = false

	go func() {
		for _ = range time.Tick(time.Minute) {
			m.Greeting = getGreeting()
		}
	}()

	views.ModelWithTemplate(m, "patient-room/patient_room.tmpl")
}

func (p *Page) OnUnload(client *ws.Client) {
	if client != nil {
		client.Unsubscribe("pong")
		client.Unsubscribe("lights1")
	}
}

func (p *Page) OnConnect(client *ws.Client) {
	client.Subscribe("pong", func(msg string) {
		println("Received pong:", msg)
		println("Latency:", time.Since(p.model.lastPing).String())
	})

	client.Subscribe("lights1", func(state int) {
		println("got lights")
		if state > 0 {
			p.model.LightOn = true
		} else {
			p.model.LightOn = false
		}
	})

	p.model.Connected = true
}

func (p *Page) OnDisconnect() {
	p.model.Connected = false
}
