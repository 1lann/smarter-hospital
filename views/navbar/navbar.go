// +build js

package navbar

import (
	"github.com/1lann/smarter-hospital/views"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/oskca/gopherjs-vue"
)

var jQuery = jquery.NewJQuery

// Item represents an item on a navbar.
type Item struct {
	*js.Object
	Name   string `js:"name"`
	Active bool   `js:"active"`
	Path   string `js:"path"`
}

// Model represents the mode for the navigation bar.
type Model struct {
	*js.Object
	Left  []Item `js:"left"`
	Right []Item `js:"right"`
}

// addPage adds a page
func (m *Model) addPage(name string, path string, right bool) {
	currentPath := js.Global.Get("location").Get("pathname").String()

	newItem := Item{Object: js.Global.Get("Object").New()}
	newItem.Name = name
	newItem.Active = currentPath == path
	newItem.Path = path

	if right {
		m.Right = append(m.Right, newItem)
		return
	}

	m.Left = append(m.Left, newItem)
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

func init() {
	templateData := views.MustAsset("navbar/navbar.tmpl")

	vue.NewComponent(func() interface{} {
		m := &Model{Object: js.Global.Get("Object").New()}
		m.Left = make([]Item, 0)
		m.Right = make([]Item, 0)

		m.addPage("Room", "/room", false)
		m.addPage("Something", "/something", false)
		m.addPage("Alerts", "/alerts", true)

		return m
	}, string(templateData)).Register("navbar")
}
