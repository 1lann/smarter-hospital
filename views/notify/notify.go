// +build js

package notify

import (
	"time"

	"github.com/1lann/smarter-hospital/notify"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
)

const serviceWorker = "/static/js/serviceWorker.js"

type Notify struct {
	client *notify.Client
	events *Events
	alerts *Alerts
}

type Notification struct {
	*js.Object
	ID         string    `js:"id"`
	Heading    string    `js:"heading"`
	SubHeading string    `js:"subHeading"`
	Icon       string    `js:"icon"`
	Link       string    `js:"link"`
	Time       time.Time `js:"time"`
}

type Alerts struct {
	*js.Object
	Events []*Notification `js:"events"`
}

type Events struct {
	*js.Object
	Events []*Notification `js:"events"`
}

type PushNotification struct {
	*js.Object
	Body               string     `js:"body"`
	Icon               string     `js:"icon"`
	RequireInteraction bool       `js:"requireInteraction"`
	Tag                string     `js:"tag"`
	OnClick            *js.Object `js:"onClick"`
}

func (c *Notify) Init() {
	js.Global.Get("Push").Get("Permission").Call("request",
		js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} { return nil }),
		js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} { return nil }),
	)

	alerts := &Alerts{Object: js.Global.Get("Object").New()}
	alerts.Events = make([]*Notification, 0)
	c.alerts = alerts

	events := &Events{Object: js.Global.Get("Object").New()}
	events.Events = make([]*Notification, 0)
	c.events = events

	views.ComponentWithTemplate(func() interface{} {
		return alerts
	}, "notify/alerts.tmpl", "mobile").Register("alerts")

	views.ComponentWithTemplate(func() interface{} {
		return events
	}, "notify/events.tmpl").Register("events")
}

func (c *Notify) OnConnect(wsClient *ws.Client) {
	client := notify.NewClient(wsClient)
	c.client = client

	client.OnNotification(func(n notify.Notification) {
		if n.Alert && !n.Dismissed {
			vue.Push(c.alerts.Get("events"), n)
			push := &PushNotification{Object: js.Global.Get("Object").New()}
			push.Body = n.SubHeading
			push.Icon = getIconURL(n.Icon)
			push.RequireInteraction = true
			push.Tag = string(n.ID)
			push.OnClick = js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
				wait := make(chan bool)

				go func() {
					client.Dismiss(string(n.ID))
					wait <- true
				}()

				js.Global.Get("window").Call("focus")
				js.Global.Get("Push").Call("close", n.ID)
				<-wait
				js.Global.Get("window").Get("location").Set("href", n.Link)
				return nil
			})
			js.Global.Get("Push").Call("create", n.Heading, push)
		}

		vue.Push(c.events.Get("events"), n)
	})

	client.OnDismiss(func(id string) {
		js.Global.Get("Push").Call("close", id)

		for i, event := range c.alerts.Events {
			if event.ID == id {
				vue.Splice(c.alerts.Get("events"), i, 1)
				return
			}
		}
	})
}

const iconLocation = "/static/icons/"

func getIconURL(icon string) string {
	switch icon {
	case "warning sign":
		return iconLocation + "warning_sign.png"
	case "doctor":
		return iconLocation + "doctor.png"
	case "plug":
		return iconLocation + "plug.png"
	case "heartbeat":
		return iconLocation + "heartbeat.png"
	default:
		return iconLocation + "alert.png"
	}
}
