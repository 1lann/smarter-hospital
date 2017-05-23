// +build js

package views

import (
	"github.com/1lann/smarter-hospital/views/charts"
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

func NewRealtimeChart(elementID string, lineColor interface{}, data []*charts.Point) *charts.Chart {
	cd := charts.NewChartData()

	d := cd.NewDataset()
	d.LineTension = 0.2

	opts := charts.NewOptions()
	opts.Responsive = true
	opts.MaintainAspectRatio = true

	ticks := charts.NewTicks()
	ticks.BeginAtZero = true

	opts.NewYAxes().Ticks = ticks

	anim := charts.NewAnimation()
	anim.Easing = "linear"
	anim.Duration = 300
	opts.Animation = anim

	d.BorderColor = lineColor
	d.BackgroundColor = "rgba(255,255,255,0)"
	d.PointRadius = 0

	axes := opts.NewXAxes()
	axes.Type = "time"
	axes.Position = "bottom"

	t := axes.NewTime()
	t.Unit = "second"
	t.UnitStepSize = 100000

	opts.NewTooltips().Enabled = false
	opts.NewLegend().Display = false

	for _, p := range data {
		d.Push(p)
	}

	return charts.NewChart(elementID, "line", cd, opts)
}

func init() {
	if js.Global.Get("location").Get("protocol").String() == "https:" {
		Address = "https://" + js.Global.Get("location").Get("host").String()
	} else {
		Address = "http://" + js.Global.Get("location").Get("host").String()
	}
}
