package charts

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

// Dataset represents a chart's dataset
type Dataset struct {
	*js.Object
	Data            []interface{} `js:"data"`
	Label           string        `js:"label"`
	BackgroundColor string        `js:"backgroundColor"`
	BorderColor     string        `js:"borderColor"`
	BorderWidth     int           `js:"borderWidth"`
	XAxisID         string        `js:"xAxisID"`
	YAxisID         string        `js:"yAxisID"`
}

// ChartData represents the chart data object to provided to render charts.
type ChartData struct {
	*js.Object
	Labels   []string   `js:"labels"`
	Datasets []*Dataset `js:"datasets"`
}

// Options represents the options to be provided to chart instance.
type Options struct {
	*js.Object
	Responsive          bool       `js:"responsive"`
	MaintainAspectRatio bool       `js:"maintainAspectRatio"`
	Scales              *Scales    `js:"scales"`
	Animation           *Animation `js:"animation"`
}

type Animation struct {
	*js.Object
	Duration int    `js:"duration"`
	Easing   string `js:"easing"`
}

// Scales represents the scales option.
type Scales struct {
	*js.Object
	yAxes []*Axes `js:"yAxes"`
	xAxes []*Axes `js:"xAxes"`
}

// Axes represents the option of a scale axes.
type Axes struct {
	*js.Object
	Type     string `js:"type"`
	Position string `js:"position"`
	Stacked  bool   `js:"stacked"`
	Display  bool   `js:"display"`
	Ticks    *Ticks `js:"ticks"`
}

// Ticks represents the options for the ticks on an axes.
type Ticks struct {
	*js.Object
	BeginAtZero   bool `js:"beginAtZero"`
	Min           int  `js:"min"`
	Max           int  `js:"max"`
	MaxTicksLimit int  `js:"maxTicksLimit"`
	FixedStepSize int  `js:"fixedStepSize"`
	StepSize      int  `js:"stepSize"`
	SuggestedMax  int  `js:"suggestedMax"`
	SuggestedMin  int  `js:"suggestedMin"`
}

// Chart represents a rendered chart.
type Chart struct {
	*js.Object
	Type    string     `js:"type"`
	Data    *ChartData `js:"data"`
	Options *Options   `js:"options"`
}

// NewDataset creates a datasets, adds it to the chart data and returns it.
func (c *ChartData) NewDataset() *Dataset {
	dataset := &Dataset{Object: js.Global.Get("Object").New()}
	c.Datasets = append(c.Datasets, dataset)
	return dataset
}

// NewChartData returns a new instance of chart data.
func NewChartData() *ChartData {
	cd := &ChartData{Object: js.Global.Get("Object").New()}
	cd.Datasets = make([]*Dataset, 0)
	return cd
}

// NewOptions returns a new options object.
func NewOptions() *Options {
	opts := &Options{Object: js.Global.Get("Object").New()}
	opts.Scales = &Scales{Object: js.Global.Get("Object").New()}
	opts.Scales.xAxes = make([]*Axes, 0)
	opts.Scales.yAxes = make([]*Axes, 0)
	return opts
}

// NewYAxes creates a new Axes object, adds it to the options as a YAxes and
// then returns it for modification.
func (o *Options) NewYAxes() *Axes {
	axes := &Axes{Object: js.Global.Get("Object").New()}
	o.Scales.yAxes = append(o.Scales.yAxes, axes)

	return axes
}

// NewXAxes creates a new Axes object, adds it to the options as a XAxes and
// then returns it for modification.
func (o *Options) NewXAxes() *Axes {
	axes := &Axes{Object: js.Global.Get("Object").New()}
	o.Scales.xAxes = append(o.Scales.xAxes, axes)

	return axes
}

// NewAnimation returns a new Animation options object.
func NewAnimation() *Animation {
	return &Animation{Object: js.Global.Get("Object").New()}
}

// NewTicks returns a new ticks option value.
func NewTicks() *Ticks {
	return &Ticks{Object: js.Global.Get("Object").New()}
}

// NewChart renders a chart and returns the resulting object.
func NewChart(elementID string, chartType string, data *ChartData,
	options *Options) *Chart {
	chartArgs := &Chart{Object: js.Global.Get("Object").New()}
	chartArgs.Type = chartType
	chartArgs.Data = data
	chartArgs.Options = options

	return &Chart{
		Object: js.Global.Get("Chart").New(dom.GetWindow().Document().
			GetElementByID("unavailable-chart").Underlying(), chartArgs),
	}
}

// Update updates the rendered chart, to be called when the underlying data
// changes.
func (c *Chart) Update() {
	c.Object.Call("update")
}
