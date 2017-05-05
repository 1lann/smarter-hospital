// +build js

package navbar

import (
	"github.com/1lann/smarter-hospital/views"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/oskca/gopherjs-vue"
)

var jQuery = jquery.NewJQuery

// Model represents the mode for the navigation bar.
type Model struct {
	*js.Object
}

// ToggleMenu toggles the navbar menu.
func (m *Model) ToggleMenu() {
	jQuery(".mobile-navbar.sidebar.menu").Call("sidebar", js.M{
		"defaultTransition": js.M{
			"computer": js.M{
				"left":   "overlay",
				"right":  "overlay",
				"top":    "overlay",
				"bottom": "overlay",
			},
			"mobile": js.M{
				"left":   "overlay",
				"right":  "overlay",
				"top":    "overlay",
				"bottom": "overlay",
			},
		},
		"onHide": js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
			jQuery(".pusher").AddClass("fullsize")
			return nil
		}),
		"onHidden": js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
			jQuery(".pusher").RemoveClass("fullsize")
			return nil
		}),
	})
	jQuery(".mobile-navbar.sidebar.menu").Call("sidebar", "toggle")
}

// CallNurse sends an alert to the nurse.
func (m *Model) CallNurse() {

}

func init() {
	templateData, err := views.Asset("patient-navbar/patient_navbar.tmpl")
	if err != nil {
		panic(err)
	}

	vue.NewComponent(func() interface{} {
		m := &Model{Object: js.Global.Get("Object").New()}
		return m
	}, string(templateData)).Register("patient-navbar")
}
