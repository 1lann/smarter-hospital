// +build js

package climate

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/1lann/smarter-hospital/logic"
	"github.com/1lann/smarter-hospital/modules/climate"
	"github.com/1lann/smarter-hospital/views"
	"github.com/1lann/smarter-hospital/views/comps"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gopherjs/gopherjs/js"
)

const (
	heatingIcon = "orange sun"
	coolingIcon = "blue fa fa-snowflake-o"
	idleIcon    = "grey asterisk"
	offIcon     = "grey power"
)

type Climate struct {
	moduleID  string
	item      *comps.Item
	component *ClimateComponent
}

type ClimateComponent struct {
	*js.Object
	On                 bool    `js:"on"`
	State              int     `js:"state"`
	CurrentTemperature float64 `js:"currentTemperature"`
	TargetTemperature  int     `js:"targetTemperature"`
	moduleID           string
}

var lightsComponent *Climate

func (c *Climate) Init(moduleID string) {
	item := &comps.Item{
		Object: js.Global.Get("Object").New(),
	}
	item.ID = moduleID
	item.Name = "Climate control"
	item.Component = "climate"
	item.Heading = "Loading..."
	item.Icon = idleIcon
	item.Available = true
	item.Active = false

	component := &ClimateComponent{
		Object:   js.Global.Get("Object").New(),
		moduleID: moduleID,
	}
	component.State = climate.StateOff
	component.CurrentTemperature = 24
	component.TargetTemperature = 24
	component.On = false

	views.ComponentWithTemplate(func() interface{} {
		return component
	}, "climate/climate.tmpl").Register("climate")

	c.moduleID = moduleID
	c.item = item
	c.component = component
}

func (c *ClimateComponent) SetTemperature(state int) {
	if state < 18 && state > 27 {
		return
	}

	go func(state int) {
		resp, err := http.Get(views.Address + "/climate/set/" + strconv.Itoa(state))
		if err != nil {
			println("set temp state:", err)
			return
		}

		defer resp.Body.Close()
	}(state)
}

func (c *ClimateComponent) Turn(state string) {
	go func(state string) {
		resp, err := http.Get(views.Address + "/climate/turn/" + state)
		if err != nil {
			println("set turn state:", err)
			return
		}

		defer resp.Body.Close()
	}(state)
}

func (c *Climate) onEvent(evt logic.ClimateState) {
	c.component.CurrentTemperature = math.Ceil((evt.CurrentTemperature*10)-0.5) / 10
	c.component.TargetTemperature = evt.TargetTemperature
	c.component.State = int(evt.State)
	c.component.On = evt.On

	if !evt.On {
		c.item.Heading = "Climate control turned off"
		c.item.Icon = offIcon
		return
	}

	switch evt.State {
	case climate.StateOff:
		c.item.Heading = "Maintaing " + strconv.Itoa(evt.TargetTemperature) + "°C"
		c.item.Icon = idleIcon
	case climate.StateCooling:
		c.item.Heading = "Cooling room to " + strconv.Itoa(evt.TargetTemperature) + "°C"
		c.item.Icon = coolingIcon
	case climate.StateHeating:
		c.item.Heading = "Heating room to " + strconv.Itoa(evt.TargetTemperature) + "°C"
		c.item.Icon = heatingIcon
	}
}

func (c *Climate) OnConnect(client *ws.Client) {
	client.Subscribe("climatecontrol", c.onEvent)

	resp, err := http.Get(views.Address + "/climate/get")
	if err != nil {
		println("climate get:", err)
		return
	}

	defer resp.Body.Close()

	var state logic.ClimateState
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&state)
	if err != nil {
		println("climate json fail")
		return
	}

	c.onEvent(state)
}

func (c *Climate) Item() *comps.Item {
	return c.item
}

func (c *Climate) OnDisconnect() {}

func (c *Climate) OnModuleConnect() {}

func (c *Climate) OnModuleDisconnect() {}
