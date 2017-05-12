package logic

import (
	"log"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/modules/lights"
	"github.com/1lann/smarter-hospital/modules/ping"
	"github.com/1lann/smarter-hospital/ws"
)

// PingLogic ...
type PingLogic struct {
	wsServer *ws.Server
	message  string
}

// Handle ...
func (p *PingLogic) Handle(s *core.Server, module ping.Module) {
	log.Println("logic handler:", module.Info())
	log.Println("last stored value:", p.message)

	p.wsServer.Emit("pong", module.Info())

	p.message = module.Info()
}

// LightLogic ...
type LightLogic struct {
	wsServer *ws.Server
}

// Handle ...
func (p *LightLogic) Handle(s *core.Server, module lights.Module) {
	p.wsServer.Emit("lights1", module.Info())
}

// Register registers logic components in the logic package using the
// provided WebSocket server as the place to emit events.
func Register(wsServer *ws.Server) {
	core.RegisterEventLogic([]string{"ping1"}, &PingLogic{wsServer: wsServer})
	core.RegisterEventLogic([]string{"lights1"}, &LightLogic{wsServer: wsServer})
}
