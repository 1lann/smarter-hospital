// +build js

package views

import (
	"github.com/oskca/gopherjs-vue"
)

// ModelWithTemplate returns a new Vue ViewModel with a model
func ModelWithTemplate(model interface{},
	templatePath string) *vue.ViewModel {
	opt := vue.NewOption()
	opt.El = "#app"
	opt.SetDataWithMethods(model)
	templateData := MustAsset(templatePath)

	opt.Template = string(templateData)
	return opt.NewViewModel()
}
