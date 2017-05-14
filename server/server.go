package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/logic"
	"github.com/1lann/smarter-hospital/store"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gin-gonic/gin"

	"github.com/1lann/smarter-hospital/views"
	_ "github.com/1lann/smarter-hospital/views/imports"

	"github.com/1lann/smarter-hospital/modules/heartrate"
	"github.com/1lann/smarter-hospital/modules/lights"
	"github.com/1lann/smarter-hospital/modules/ultrasonic"
)

var server *core.Server

var webPath = os.Getenv("GOPATH") + "/src/github.com/1lann/smarter-hospital/server"

var connectedModules = make(map[string]bool)
var connectedMutex = new(sync.Mutex)

func main() {
	err := store.Connect(store.ConnectOpts{
		Address:  "127.0.0.1:27017",
		Database: "smarter-hospital",
	})
	if err != nil {
		panic(err)
	}

	core.SetupModule("lights", "light1", lights.Settings{
		Pin:               11,
		AnimationDuration: time.Second,
	})

	core.SetupModule("ultrasonic", "ultrasonic1", ultrasonic.Settings{
		TriggerPin:       5,
		EchoPin:          6,
		ContactThreshold: 2,
	})

	core.SetupModule("heartrate", "heartrate1", heartrate.Settings{
		PeakThreshold: 530,
		Pin:           0,
	})

	wsServer := ws.NewServer()

	core.RegisterConnect(func(moduleID string) {
		connectedMutex.Lock()
		connectedModules[moduleID] = true
		connectedMutex.Unlock()
		wsServer.Emit("moduleConnect", moduleID)
	})

	core.RegisterDisconnect(func(moduleID string) {
		connectedMutex.Lock()
		delete(connectedModules, moduleID)
		connectedMutex.Unlock()
		wsServer.Emit("moduleDisconnect", moduleID)
	})

	core.RegisterEventHandler(func(moduleID string, event interface{}) {
		wsServer.Emit(moduleID, event)
	})

	logic.Register(wsServer)

	server, err = core.NewServer("0.0.0.0:5000")
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.LoadHTMLFiles(webPath + "/view.tmpl")

	allPages := views.AllPages()
	for pagePath, page := range allPages {
		r.GET(pagePath, pageHandler(page.Title()))
	}

	r.GET("/ws", func(c *gin.Context) {
		wsServer.Handle(c.Request, c.Writer)
	})

	// TODO: implement
	r.GET("/module/info/:moduleid", handleInfo)
	r.POST("/module/action/:moduleid", handleAction)
	r.GET("/connected-modules", func(c *gin.Context) {
		var results []string
		connectedMutex.Lock()
		for module := range connectedModules {
			results = append(results, module)
		}
		connectedMutex.Unlock()
		c.JSON(http.StatusOK, results)
	})

	r.Static("/static", webPath+"/vendor")

	log.Println("Server is running!")
	r.Run()
}

func pageHandler(title string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "view.tmpl", gin.H{
			"Title": title,
			"User": gin.H{
				"firstName": "John",
				"lastName":  "Smith",
			},
		})
	}
}
