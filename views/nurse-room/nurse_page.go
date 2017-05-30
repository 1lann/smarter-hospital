// +build js

package nurseroom

import (
	"strings"

	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/views/climate"
	"github.com/1lann/smarter-hospital/views/comps"
	"github.com/1lann/smarter-hospital/views/contact"
	"github.com/1lann/smarter-hospital/views/heartrate"
	"github.com/1lann/smarter-hospital/views/lights"
	"github.com/1lann/smarter-hospital/views/notify"
	_ "github.com/1lann/smarter-hospital/views/nurse-navbar"
	"github.com/1lann/smarter-hospital/views/occupancy"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
	vue "github.com/oskca/gopherjs-vue"
)

var modules = map[string]comps.Component{
	"ultrasonic1":    &contact.Contact{},
	"heartrate1":     &heartrate.HeartRate{},
	"lights1":        &lights.Lights{},
	"climatecontrol": &climate.Climate{},
	"proximity1":     &occupancy.Occupancy{},
}

var notifier *notify.Notify

func (m *Model) SelectComponent(component string) {
	for _, category := range m.Categories {
		for _, item := range category.Items {
			if item.Component == component {
				item.Active = true
				js.Global.Get("location").Set("hash", "#"+item.ID)
			} else {
				item.Active = false
			}
		}
	}

	m.ShowMenu = false
	m.ViewComponent = component
}

func (m *Model) DisplayMenu() {
	for _, category := range m.Categories {
		for _, item := range category.Items {
			item.Active = false
		}
	}

	m.ShowMenu = true
	js.Global.Get("location").Set("hash", "")
}

func populateCategories(m *Model) {
	// eventsCat := &comps.Category{Object: js.Global.Get("Object").New()}
	//
	// eventsCat.Heading = ""
	// eventsCat.SubHeading = ""
	// eventsCat.Icon = ""
	// eventsCat.Items = make([]*comps.Item, 0)

	// eventItem := &comps.Item{Object: js.Global.Get("Object").New()}
	// eventItem.ID = "events"
	// eventItem.Name = "Event log"
	// eventItem.Icon = "history"
	// eventItem.Heading = "Events"
	// eventItem.SubHeading = "Show event log"
	// eventItem.Component = "events"
	// eventItem.Available = true
	// eventItem.Active = false
	// eventsCat.Items = append(eventsCat.Items, eventItem)

	// m.Categories = append(m.Categories, eventsCat)

	roomControls := &comps.Category{Object: js.Global.Get("Object").New()}

	roomControls.Heading = "Room controls"
	roomControls.SubHeading = ""
	roomControls.Icon = "settings"
	roomControls.Items = make([]*comps.Item, 0)

	roomControls.Items = append(roomControls.Items, modules["lights1"].Item())
	roomControls.Items = append(roomControls.Items, modules["climatecontrol"].Item())
	m.Categories = append(m.Categories, roomControls)

	healthCategory := &comps.Category{
		Object: js.Global.Get("Object").New(),
	}

	healthCategory.Heading = "Patient health"
	healthCategory.SubHeading = ""
	healthCategory.Icon = "green plus"
	healthCategory.Items = make([]*comps.Item, 0)

	healthCategory.Items = append(healthCategory.Items, modules["ultrasonic1"].Item())
	healthCategory.Items = append(healthCategory.Items, modules["heartrate1"].Item())
	healthCategory.Items = append(healthCategory.Items, modules["proximity1"].Item())
	m.Categories = append(m.Categories, healthCategory)
}

func (p *Page) OnLoad() {
	pageModel := &Model{
		Object: js.Global.Get("Object").New(),
	}

	type Unavailable struct {
		*js.Object
	}

	u := &Unavailable{
		Object: js.Global.Get("Object").New(),
	}

	uClosure := func() interface{} {
		return u
	}

	opt := vue.NewOption()
	opt.Data = uClosure
	opt.Template = string(views.MustAsset("nurse-room/unavailable.tmpl"))
	opt.AddProp("connected")
	opt.OnLifeCycleEvent(vue.EvtBeforeCreate, func(vm *vue.ViewModel) {
		vm.Options.Set("methods", js.MakeWrapper(u))
	})

	{
		//
		// charts.SetGlobalFontFamily("Lato")
		// charts.SetGlobalFontSize(14)
		//
		// shouldContinue := true
		//
		// opt.OnLifeCycleEvent(vue.EvtMounted, func(vm *vue.ViewModel) {
		// 	println("Created")
		//
		// 	grad := charts.NewGradient("unavailable-chart", 0, 500)
		// 	grad.AddColorStop(0, charts.Color(0x4CAF50))
		// 	grad.AddColorStop(1, charts.Color(0xF44336))
		//
		// 	data := make([]*charts.Point, 0)
		// 	for i := -100; i < 0; i++ {
		// 		data = append(data, charts.NewPoint(float64(time.Now().Add(
		// 			time.Millisecond*time.Duration(i)*200).UnixNano())/1000000, 1))
		// 	}
		//
		// 	chart := views.NewRealtimeChart("unavailable-chart", grad, data)
		//
		// 	go func() {
		// 		r := rand.New(rand.NewSource(0))
		//
		// 		for shouldContinue {
		// 			num := r.Int()%10 + 1
		// 			chart.Data.Datasets[0].PushAndShift(charts.NewPoint(float64(time.Now().UnixNano())/1000000, num))
		// 			chart.Options.Scales.XAxes[0].Time.Min = float64(time.Now().Add(time.Millisecond*-19500).UnixNano()) / 1000000
		// 			chart.Options.Scales.XAxes[0].Time.Max = float64(time.Now().Add(time.Millisecond*-500).UnixNano()) / 1000000
		// 			chart.Update()
		// 			time.Sleep(time.Millisecond * 200)
		// 		}
		//
		// 		println("Exited goroutine")
		//
		// 	}()
		// })
		// opt.OnLifeCycleEvent(vue.EvtDestroyed, func(vm *vue.ViewModel) {
		// 	println("Destroyed")
		// 	shouldContinue = false
		//
		// })
	}

	opt.NewComponent().Register(comps.UnavailableView)

	for moduleID, module := range modules {
		module.Init(moduleID)
	}

	notifier = &notify.Notify{}
	notifier.Init()

	pageModel.Categories = make([]*comps.Category, 0)
	populateCategories(pageModel)

	p.model = pageModel

	pageModel.Name = views.GetUser().FirstName + " " + views.GetUser().LastName
	pageModel.Connected = false

	hash := strings.TrimPrefix(js.Global.Get("location").Get("hash").String(), "#")

	pageModel.Mobile = js.Global.Get("window").Get("innerWidth").Int() <= 700
	pageModel.ShowMenu = true

	if hash != "" {
		p.showHash(hash)
	} else if !pageModel.Mobile {
		pageModel.ViewComponent = pageModel.Categories[0].Items[0].Component
		pageModel.Categories[0].Items[0].Active = true
		js.Global.Get("location").Set("hash", "#"+pageModel.Categories[0].Items[0].ID)
	} else {
		pageModel.ViewComponent = ""
	}

	js.Global.Get("window").Set("onhashchange", js.MakeFunc(func(this *js.Object, arguments []*js.Object) interface{} {
		p.showHash(strings.TrimPrefix(js.Global.Get("location").Get("hash").String(), "#"))
		return nil
	}))

	views.ModelWithTemplate(pageModel, "nurse-room/nurse_room.tmpl")
}

func (p *Page) showHash(hash string) {
	for _, category := range p.model.Categories {
		for _, item := range category.Items {
			if item.ID == hash {
				item.Active = true
				p.model.ShowMenu = false
				p.model.ViewComponent = item.Component
			} else {
				item.Active = false
			}
		}
	}
}

func (p *Page) OnUnload(client *ws.Client) {
	// TODO: Consider the need for this
}

func (p *Page) OnConnect(client *ws.Client) {
	if p.connected {
		js.Global.Get("location").Call("reload")
	}
	p.connected = true
	p.model.Connected = true

	client.Subscribe("moduleConnect", func(moduleID string) {
		module, found := modules[moduleID]
		if !found {
			return
		}

		module.Item().Available = true
		if module.Item().Active {
			p.model.ViewComponent = module.Item().Component
		}

		module.OnModuleConnect()
	})

	client.Subscribe("moduleDisconnect", func(moduleID string) {
		module, found := modules[moduleID]
		if !found {
			return
		}

		module.Item().Available = false
		if module.Item().Active {
			p.model.ViewComponent = comps.UnavailableView
		}

		module.OnModuleDisconnect()
	})

	notifier.OnConnect(client)

	connectedModules, err := views.ConnectedModules()
	if err != nil {
		println("connected modules error:", err.Error())
	}

	for moduleID, module := range modules {
		module.OnConnect(client)

		if moduleID[len(moduleID)-1] > '0' && moduleID[len(moduleID)-1] < '9' {
			if isInList(moduleID, connectedModules) {
				module.Item().Available = true
				if module.Item().Active {
					p.model.ViewComponent = module.Item().Component
				}
				module.OnModuleConnect()
			} else {
				module.Item().Available = false
				if module.Item().Active {
					p.model.ViewComponent = comps.UnavailableView
				}

				module.OnModuleDisconnect()
			}
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
