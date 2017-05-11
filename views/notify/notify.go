// +build js

package notify

import (
	"github.com/1lann/smarter-hospital/views"
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
)

// Notification represents a notification.
type Notification struct {
	*js.Object
	Message  string `js:"message"`
	Severity string `js:"severity"`
	Page     string `js:"page"`
}

// Model represents the model for the notifications view.
type Model struct {
	*js.Object
	Notifications []Notification `js:"notifications"`
}

func init() {
	templateData := views.MustAsset("notify/notify.tmpl")

	vue.NewComponent(func() interface{} {
		return &Model{Object: js.Global.Get("Object").New()}
	}, string(templateData)).Register("notify")
}
