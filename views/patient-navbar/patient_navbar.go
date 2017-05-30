// +build js

package navbar

import (
	"net/http"
	"time"

	"github.com/1lann/smarter-hospital/views"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jQuery = jquery.NewJQuery

// Model represents the mode for the navigation bar.
type Model struct {
	*js.Object
	Name       string `js:"name"`
	Date       string `js:"date"`
	Time       string `js:"time"`
	RoomNumber string `js:"roomNumber"`
	Connected  bool   `js:"connected"`
	Nurse      *Nurse `js:"nurse"`
}

type Nurse struct {
	*js.Object
	Status     string `js:"status"`
	Icon       string `js:"icon"`
	AllowCalls bool   `js:"allowCalls"`
	Color      string `js:"color"`
}

// CallNurse sends an alert to the nurse.
func (m *Model) CallNurse() {
	denied := make(chan bool)

	jquery.NewJQuery(".ui.modal.nurse-modal .progress .bar").SetCss("width", "0%")
	jquery.NewJQuery(".ui.modal.nurse-modal").Call("modal", js.M{
		"inverted": true,
		"closable": false,
		"onVisible": js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
			jquery.NewJQuery(".ui.modal.nurse-modal .progress .bar").SetCss("width", "100%")

			go func() {
				select {
				case <-time.After(time.Second * 5):
					// TODO: Make request
					jquery.NewJQuery(".ui.modal.nurse-modal").Call("modal", "hide")

					go func() {
						resp, err := http.Get(views.Address + "/notify/call")
						if err != nil {
							println("call nurse:", err)
							return
						}

						defer resp.Body.Close()
					}()

					m.Nurse.AllowCalls = false
					m.Nurse.Icon = "checkmark"
					m.Nurse.Color = "green"
					m.Nurse.Status = "Nurse called"

					time.Sleep(time.Second * 10)

					m.Nurse.AllowCalls = true
					m.Nurse.Icon = "doctor"
					m.Nurse.Color = "red"
					m.Nurse.Status = "Call nurse"
				case <-denied:
					println("cancelled")
				}
			}()

			return nil
		}),
		"onDeny": js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
			denied <- true
			return nil
		}),
	}).Call("modal", "show")
}

func init() {
	views.ComponentWithTemplate(func() interface{} {
		m := &Model{Object: js.Global.Get("Object").New()}
		m.RoomNumber = "025"
		m.Name = views.GetUser().FirstName + " " + views.GetUser().LastName

		m.Date = time.Now().Format("Monday, _2 Jan 2006")
		m.Time = time.Now().Format("3:04:05 PM")

		m.Nurse = &Nurse{Object: js.Global.Get("Object").New()}

		m.Nurse.Status = "Call nurse"
		m.Nurse.Icon = "doctor"
		m.Nurse.AllowCalls = true
		m.Nurse.Color = "red"

		go func() {
			for _ = range time.Tick(time.Second) {
				m.Date = time.Now().Format("Monday, _2 Jan 2006")
				m.Time = time.Now().Format("3:04:05 PM")
			}
		}()

		return m
	}, "patient-navbar/patient_navbar.tmpl", "connected").
		Register("patient-navbar")
}
