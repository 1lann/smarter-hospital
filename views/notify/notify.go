// +build js

package notify

import (
	"time"

	"github.com/1lann/smarter-hospital/notify"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/dustin/go-humanize"
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
	ID           string `js:"id"`
	Heading      string `js:"heading"`
	SubHeading   string `js:"subHeading"`
	Location     string `js:"location"`
	Icon         string `js:"icon"`
	Link         string `js:"link"`
	Time         string `js:"time"`
	InternalTime int64  `js:"internalTime"`
}

type Alerts struct {
	*js.Object
	Events   []*Notification `js:"events"`
	ShowMenu bool            `js:"showMenu"`
	client   *notify.Client
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

func (a *Alerts) Click(id string) {
	wait := make(chan bool)

	link := ""
	for _, alert := range a.Events {
		if alert.ID == id {
			link = alert.Link
		}
	}

	go func() {
		a.client.Dismiss(id)
		wait <- true
	}()

	js.Global.Get("Push").Call("close", id)
	<-wait
	js.Global.Get("window").Get("location").Set("href", link)
}

func (c *Notify) Init() {
	js.Global.Get("Push").Get("Permission").Call("request",
		js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} { return nil }),
		js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} { return nil }),
	)

	alerts := &Alerts{Object: js.Global.Get("Object").New()}
	alerts.Events = make([]*Notification, 0)
	alerts.ShowMenu = false
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

	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			for _, event := range events.Events {
				event.Time = humanize.Time(time.Unix(event.InternalTime, 0))
			}
		}
	}()
}

func (c *Notify) OnConnect(wsClient *ws.Client) {
	client := notify.NewClient(wsClient)
	c.client = client
	c.alerts.client = client

	client.OnNotification(func(n notify.Notification) {
		notif := &Notification{Object: js.Global.Get("Object").New()}
		notif.ID = n.ID.Hex()
		notif.Heading = n.Heading
		notif.SubHeading = n.SubHeading
		notif.Location = n.Location
		notif.Icon = n.Icon
		notif.Link = n.Link
		notif.InternalTime = n.Time.Unix()
		notif.Time = humanize.Time(n.Time)

		if n.Alert && !n.Dismissed {
			vue.Push(c.alerts.Get("events"), notif)
			push := &PushNotification{Object: js.Global.Get("Object").New()}
			push.Body = n.SubHeading
			push.Icon = getIconURL(n.Icon)
			push.RequireInteraction = true
			push.Tag = n.ID.Hex()
			push.OnClick = js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
				wait := make(chan bool)

				go func() {
					client.Dismiss(n.ID.Hex())
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

		vue.Push(c.events.Get("events"), notif)
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

	notifications, err := client.Start()
	if err != nil {
		return
	}

	for _, notification := range notifications {
		n := &Notification{Object: js.Global.Get("Object").New()}
		n.ID = notification.ID.Hex()
		n.Heading = notification.Heading
		n.SubHeading = notification.SubHeading
		n.Location = notification.Location
		n.Icon = notification.Icon
		n.Link = notification.Link
		n.InternalTime = notification.Time.Unix()
		n.Time = humanize.Time(notification.Time)

		if notification.Alert {
			vue.Push(c.alerts.Get("events"), n)

			push := &PushNotification{Object: js.Global.Get("Object").New()}
			push.Body = n.SubHeading
			push.Icon = getIconURL(n.Icon)
			push.RequireInteraction = true
			push.Tag = string(n.ID)
			push.OnClick = js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
				go func() {
					wait := make(chan bool)

					go func() {
						client.Dismiss(n.ID)
						wait <- true
					}()

					js.Global.Get("window").Call("focus")
					js.Global.Get("Push").Call("close", n.ID)
					<-wait
					js.Global.Get("window").Get("location").Set("href", n.Link)
				}()
				return nil
			})
			js.Global.Get("Push").Call("create", n.Heading, push)
		}

		vue.Push(c.events.Get("events"), n)
	}
}

const iconLocation = "/static/icons/"

func getIconURL(icon string) string {
	switch icon {
	case "red warning sign":
		return iconLocation + "warning_sign.png"
	case "red doctor":
		return iconLocation + "doctor.png"
	case "red plug":
		return iconLocation + "plug.png"
	case "red heartbeat":
		return iconLocation + "heartbeat.png"
	default:
		return iconLocation + "alert.png"
	}
}
