// +build js

package patientroom

import (
	"math/rand"
	"time"

	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/views/comps"
	"github.com/1lann/smarter-hospital/views/contact"
	"github.com/1lann/smarter-hospital/views/heartrate"
	"github.com/1lann/smarter-hospital/views/lights"
	_ "github.com/1lann/smarter-hospital/views/patient-navbar"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
	vue "github.com/oskca/gopherjs-vue"

	"github.com/1lann/smarter-hospital/views/charts"
)

var modules = map[string]comps.Component{
	"ultrasonic1": &contact.Contact{},
	"heartrate1":  &heartrate.HeartRate{},
	"lights1":     &lights.Lights{},
}

func (m *Model) SelectComponent(component string) {
	for _, category := range m.Categories {
		for _, item := range category.Items {
			if item.Component == component {
				item.Active = true
			} else {
				item.Active = false
			}
		}
	}

	m.ViewComponent = component
}

func getGreeting() string {
	hour := time.Now().Hour()
	if hour < 5 {
		return "evening"
	} else if hour < 12 {
		return "morning"
	} else if hour < 18 {
		return "afternoon"
	}

	return "evening"
}

func populateCategories(m *Model) {
	roomControls := &comps.Category{
		Object: js.Global.Get("Object").New(),
	}

	roomControls.Heading = "Room controls"
	roomControls.SubHeading = ""
	roomControls.Icon = "settings"
	roomControls.Items = make([]*comps.Item, 0)

	roomControls.Items = append(roomControls.Items, modules["lights1"].Item())
	m.Categories = append(m.Categories, roomControls)

	healthCategory := &comps.Category{
		Object: js.Global.Get("Object").New(),
	}

	healthCategory.Heading = "Your health"
	healthCategory.SubHeading = ""
	healthCategory.Icon = "green plus"
	healthCategory.Items = make([]*comps.Item, 0)

	healthCategory.Items = append(healthCategory.Items, modules["ultrasonic1"].Item())
	healthCategory.Items = append(healthCategory.Items, modules["heartrate1"].Item())
	m.Categories = append(m.Categories, healthCategory)
}

func (p *Page) OnLoad() {
	pageModel := &Model{
		Object: js.Global.Get("Object").New(),
	}

	type Unavailable struct {
		*js.Object
		ChartData *charts.ChartData `js:"chartData"`
	}

	u := &Unavailable{
		Object: js.Global.Get("Object").New(),
	}

	uClosure := func() interface{} {
		return u
	}

	opt := vue.NewOption()
	opt.Data = uClosure
	opt.Template = string(views.MustAsset("patient-room/unavailable.tmpl"))
	opt.AddProp("connected")
	opt.OnLifeCycleEvent(vue.EvtBeforeCreate, func(vm *vue.ViewModel) {
		vm.Options.Set("methods", js.MakeWrapper(u))
	})

	shouldContinue := true

	opt.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
		println("Created")
		cd := charts.NewChartData()
		d := cd.NewDataset()
		d.Data = []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
		d.Label = "Testing 123"
		cd.Labels = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}
		opts := charts.NewOptions()
		ticks := charts.NewTicks()
		ticks.BeginAtZero = true
		opts.NewYAxes().Ticks = ticks
		anim := charts.NewAnimation()
		anim.Easing = "linear"
		anim.Duration = 1500
		opts.Animation = anim
		shouldContinue = true

		chart := charts.NewChart("unavailable-chart", "line", cd, opts)

		go func() {
			r := rand.New(rand.NewSource(0))

			for shouldContinue {
				chart.Data.Datasets[0].Object.Get("data").Call("push", r.Int()%10)
				chart.Data.Datasets[0].Object.Get("data").Call("shift")
				chart.Update()
				time.Sleep(time.Second)
			}

			println("Exited goroutine")

		}()
	})
	opt.OnLifeCycleEvent(vue.EvtDestroyed, func(vm *vue.ViewModel) {
		println("Destroyed")
		shouldContinue = false

	})
	opt.NewComponent().Register(comps.UnavailableView)

	for moduleID, module := range modules {
		module.Init(moduleID)
	}

	pageModel.Categories = make([]*comps.Category, 0)
	populateCategories(pageModel)

	p.model = pageModel

	pageModel.Name = views.GetUser().FirstName + " " + views.GetUser().LastName
	pageModel.Greeting = getGreeting()
	pageModel.PingText = ""
	pageModel.LightOn = false
	pageModel.Connected = false
	pageModel.ViewComponent = comps.UnavailableView

	go func() {
		for _ = range time.Tick(time.Minute) {
			pageModel.Greeting = getGreeting()
		}
	}()

	views.ModelWithTemplate(pageModel, "patient-room/patient_room.tmpl")
}

func (p *Page) OnUnload(client *ws.Client) {
	// TODO: Consider the need for this
}

// TODO: Needs a lot of cleaning!

func (p *Page) OnConnect(client *ws.Client) {
	if p.connected {
		js.Global.Get("location").Call("reload")
	}
	p.connected = true
	p.model.Connected = true

	// TODO: Automate unavailability

	client.Subscribe("moduleConnect", func(moduleID string) {
		module, found := modules[moduleID]
		if !found {
			return
		}

		module.OnModuleConnect()
	})

	client.Subscribe("moduleDisconnect", func(moduleID string) {
		module, found := modules[moduleID]
		if !found {
			return
		}

		module.OnModuleDisconnect()
	})

	connectedModules, err := views.ConnectedModules()
	if err != nil {
		println("connected modules error:", err.Error())
	}

	for moduleID, module := range modules {
		module.OnConnect(client)

		if isInList(moduleID, connectedModules) {
			module.OnModuleConnect()
		} else {
			module.OnModuleDisconnect()
		}
	}
}

func isInList(item string, list []string) bool {
	for _, listItem := range list {
		if item == listItem {
			return true
		}
	}
	return false
}

func (p *Page) OnDisconnect() {
	p.model.Connected = false

	for _, module := range modules {
		module.OnDisconnect()
	}
}
