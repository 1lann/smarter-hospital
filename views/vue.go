// +build js

package views

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
)

// ModelWithTemplate returns a new Vue ViewModel with a model
func ModelWithTemplate(model interface{},
	templatePath string) *vue.ViewModel {
	opt := vue.NewOption()
	opt.El = "#app"
	opt.SetDataWithMethods(model)
	opt.Template = string(MustAsset(templatePath))
	return opt.NewViewModel()
}

// ComponentWithTemplate creates a new component with the provided initializer
// function, template path and properties.
func ComponentWithTemplate(initializer func() (model interface{}),
	templatePath string, props ...string) *vue.Component {
	model := initializer()
	modelClosure := func() interface{} {
		return model
	}

	opt := vue.NewOption()
	opt.Data = modelClosure
	opt.Template = string(MustAsset(templatePath))
	opt.AddProp(props...)
	opt.OnLifeCycleEvent(vue.EvtBeforeCreate, func(vm *vue.ViewModel) {
		vm.Options.Set("methods", js.MakeWrapper(model))
	})
	return opt.NewComponent()
}
