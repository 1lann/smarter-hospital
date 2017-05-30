package logic

import (
	"log"
	"net/http"
	"strconv"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/modules/climate"
	"github.com/1lann/smarter-hospital/notify"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gin-gonic/gin"
)

// ClimateControl ...
type ClimateControl struct {
	wsServer     *ws.Server
	notifyServer *notify.Server
	state        ClimateState
}

type ClimateState struct {
	On                 bool
	State              climate.State
	TargetTemperature  int
	CurrentTemperature float64
}

func (c *ClimateControl) GetState(g *gin.Context) {
	g.JSON(http.StatusOK, c.state)
}

func (c *ClimateControl) Turn(g *gin.Context) {
	if g.Param("onoff") == "on" {
		c.state.On = true
	} else {
		c.state.On = false
	}

	c.wsServer.Emit("climatecontrol", c.state)
	g.String(http.StatusOK, "")
}

func (c *ClimateControl) SetTemperature(g *gin.Context) {
	num, err := strconv.Atoi(g.Param("temp"))
	if err != nil {
		log.Println("climate control:", err)
		g.String(http.StatusInternalServerError, "fuck")
		return
	}

	if num >= 18 && num <= 27 {
		c.state.TargetTemperature = num
	}

	c.wsServer.Emit("climatecontrol", c.state)

	g.String(http.StatusOK, "")
}

type Warner struct {
	notifyServer      *notify.Server
	hasHeartrateAlert bool
	hasBedAlert       bool
}

type SmartLighting struct {
	hasChanged bool
}

// TODO: Occupancy if you have time

// Register registers logic components in the logic package using the
// provided WebSocket server as the place to emit events.
func Register(r *gin.Engine, wsServer *ws.Server, notifyServer *notify.Server) {
	c := &ClimateControl{
		wsServer:     wsServer,
		notifyServer: notifyServer,
		state: ClimateState{
			State:              climate.StateOff,
			TargetTemperature:  24,
			CurrentTemperature: 24,
		},
	}
	core.RegisterEventLogic([]string{"climate1", "thermometer1"}, c)

	w := &Warner{
		notifyServer:      notifyServer,
		hasHeartrateAlert: true,
		hasBedAlert:       true,
	}
	core.RegisterEventLogic([]string{"heartrate1", "ultrasonic1"}, w)

	sl := &SmartLighting{}
	core.RegisterEventLogic([]string{"proximity1"}, sl)

	r.GET("/climate/get", c.GetState)
	r.GET("/climate/set/:temp", c.SetTemperature)
	r.GET("/climate/turn/:onoff", c.Turn)
}
