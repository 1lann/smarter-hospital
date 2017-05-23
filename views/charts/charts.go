package charts

import (
	"image/color"
	"strconv"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

// Dataset represents a chart's dataset
type Dataset struct {
	*js.Object
	Label string `js:"label"`

	Data []*Point `js:"data"`

	BackgroundColor interface{} `js:"backgroundColor"`

	HoverBackgroundColor interface{} `js:"hoverBackgroundColor"`
	BorderColor          interface{} `js:"borderColor"`

	HoverBorderColor interface{} `js:"hoverBorderColor"`
	BorderWidth      int         `js:"borderWidth"`

	PointBorderColor          interface{} `js:"pointBorderColor"`
	PointBackgroundColor      interface{} `js:"pointBackgroundColor"`
	PointBorderWidth          interface{} `js:"pointBorderWidth"`
	PointRadius               int         `js:"pointRadius"`
	PointHoverRadius          int         `js:"pointHoverRadius"`
	PointHitRadius            int         `js:"pointHitRadius"`
	PointHoverBackgroundColor interface{} `js:"pointHoverBackgroundColor"`
	PointHoverBorderColor     interface{} `js:"pointHoverBorderColor"`
	PointHoverBorderWidth     interface{} `js:"pointHoverBorderWidth"`
	PointStyle                string      `js:"pointStyle"`

	LineTension float64 `js:"lineTension"`

	XAxisID string `js:"xAxisID"`
	YAxisID string `js:"yAxisID"`
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
	Tooltips            *Tooltips  `js:"tooltips"`
	Legend              *Legend    `js:"legend"`
}

// Animation represents the rendering animation options.
type Animation struct {
	*js.Object
	Duration int    `js:"duration"`
	Easing   string `js:"easing"`
}

// Scales represents the scales option.
type Scales struct {
	*js.Object
	YAxes []*Axes `js:"yAxes"`
	XAxes []*Axes `js:"xAxes"`
}

// Axes represents the option of a scale axes.
type Axes struct {
	*js.Object
	Type     string `js:"type"`
	Position string `js:"position"`
	Stacked  bool   `js:"stacked"`
	Display  bool   `js:"display"`
	Ticks    *Ticks `js:"ticks"`
	Time     *Time  `js:"time"`
}

// Tooltips represents the options for tooltips.
type Tooltips struct {
	*js.Object
	Enabled bool `js:"enabled"`
}

// Legend represents the options for the chart legend.
type Legend struct {
	*js.Object
	Display bool `js:"display"`
}

// Ticks represents the options for the ticks on an axes.
type Ticks struct {
	*js.Object
	BeginAtZero   bool    `js:"beginAtZero"`
	Min           float64 `js:"min"`
	Max           float64 `js:"max"`
	MaxTicksLimit int     `js:"maxTicksLimit"`
	FixedStepSize int     `js:"fixedStepSize"`
	StepSize      int     `js:"stepSize"`
	SuggestedMax  int     `js:"suggestedMax"`
	SuggestedMin  int     `js:"suggestedMin"`
}

// Chart represents a rendered chart.
type Chart struct {
	*js.Object
	Type    string     `js:"type"`
	Data    *ChartData `js:"data"`
	Options *Options   `js:"options"`
}

// Gradient represents a color gradient for use in charts.
type Gradient struct {
	*js.Object
}

// Time represents the option for time type charts.
type Time struct {
	*js.Object
	DisplayFormats *DisplayFormats `js:"displayFormats"`
	Max            float64         `js:"max"`
	Min            float64         `js:"min"`
	Round          string          `js:"round"`
	TooltipFormat  string          `js:"tooltipFormat"`
	Unit           string          `js:"unit"`
	UnitStepSize   int             `js:"unitStepSize"`
	MinUnit        string          `js:"minUnit"`
}

// DisplayFormats represents the time display formats.
type DisplayFormats struct {
	*js.Object
	Millisecond string `js:"millisecond"`
	Second      string `js:"second"`
	Minute      string `js:"minute"`
	Hour        string `js:"hour"`
	Day         string `js:"day"`
	Week        string `js:"week"`
	Month       string `js:"month"`
	Quarter     string `js:"quarter"`
	Year        string `js:"year"`
}

// Point represents a data point.
type Point struct {
	*js.Object
	X interface{} `js:"x"`
	Y interface{} `js:"y"`
}

// NewPoint creates a new point and returns it.
func NewPoint(x interface{}, y interface{}) *Point {
	p := &Point{Object: js.Global.Get("Object").New()}
	p.X = x
	p.Y = y
	return p
}

// NewDataset creates a datasets, adds it to the chart data and returns it.
func (c *ChartData) NewDataset() *Dataset {
	dataset := &Dataset{Object: js.Global.Get("Object").New()}
	dataset.Data = make([]*Point, 0)
	c.Datasets = append(c.Datasets, dataset)
	return dataset
}

// Push adds the provided value to data array.
func (d *Dataset) Push(p *Point) {
	d.Object.Get("data").Call("push", p)
}

// PushAndShift adds the provided value to data array, and then shifts it.
func (d *Dataset) PushAndShift(p *Point) {
	d.Object.Get("data").Call("push", p)
	d.Object.Get("data").Call("shift")
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
	opts.Scales.XAxes = make([]*Axes, 0)
	opts.Scales.YAxes = make([]*Axes, 0)
	return opts
}

// NewYAxes creates a new Axes object, adds it to the options as a YAxes and
// then returns it for modification.
func (o *Options) NewYAxes() *Axes {
	axes := &Axes{Object: js.Global.Get("Object").New()}
	o.Scales.YAxes = append(o.Scales.YAxes, axes)

	return axes
}

// NewXAxes creates a new Axes object, adds it to the options as a XAxes and
// then returns it for modification.
func (o *Options) NewXAxes() *Axes {
	axes := &Axes{Object: js.Global.Get("Object").New()}
	o.Scales.XAxes = append(o.Scales.XAxes, axes)

	return axes
}

// NewTooltips creates a new Tooltips object and returns it.
func (o *Options) NewTooltips() *Tooltips {
	o.Tooltips = &Tooltips{Object: js.Global.Get("Object").New()}
	return o.Tooltips
}

// NewLegend creates a new Legend object and returns it.
func (o *Options) NewLegend() *Legend {
	o.Legend = &Legend{Object: js.Global.Get("Object").New()}
	return o.Legend
}

// NewTime returns a new Time options object for the given axes option.
func (a *Axes) NewTime() *Time {
	a.Time = &Time{Object: js.Global.Get("Object").New()}
	return a.Time
}

// NewDisplayFormats returns a new DisplayFormats object for the time option.
func (t *Time) NewDisplayFormats() *DisplayFormats {
	t.DisplayFormats = &DisplayFormats{Object: js.Global.Get("Object").New()}
	return t.DisplayFormats
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

// NewGradient creates a new gradient for the provided canvas element and
// returns it.
func NewGradient(elementID string, width int, height int) *Gradient {
	return &Gradient{
		Object: dom.GetWindow().Document().GetElementByID(
			elementID).(*dom.HTMLCanvasElement).GetContext2d().
			Call("createLinearGradient", 0, 0, width, height),
	}
}

// AddColorStop adds a color stop to the gradient.
func (grad *Gradient) AddColorStop(offset float64, col color.Color) {
	r, g, b, a := col.RGBA()
	rs := strconv.Itoa(int(r >> (1 << 3)))
	gs := strconv.Itoa(int(g >> (1 << 3)))
	bs := strconv.Itoa(int(b >> (1 << 3)))
	as := strconv.FormatFloat(float64(a)/0xffff, 'f', 5, 64)

	grad.Call("addColorStop", offset, "rgba("+rs+","+gs+","+bs+","+as+")")
}

// Color returns a color with the given hex code.
func Color(hex uint32) color.Color {
	return color.RGBA{
		R: uint8(hex >> (1 << 4)),
		G: uint8((hex & 0xff00) >> (1 << 3)),
		B: uint8(hex & 0xff),
		A: 255,
	}
}

// SetGlobalFontFamily sets the global default font family.
func SetGlobalFontFamily(fontFamily string) {
	js.Global.Get("Chart").Get("defaults").Get("global").Set("defaultFontFamily", fontFamily)
}

// SetGlobalFontSize sets the global default font size (in pixels).
func SetGlobalFontSize(fontSize int) {
	js.Global.Get("Chart").Get("defaults").Get("global").Set("defaultFontSize", fontSize)
}
