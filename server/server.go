package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/logic"
	"github.com/1lann/smarter-hospital/notify"
	"github.com/1lann/smarter-hospital/store"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/1lann/smarter-hospital/views"
	_ "github.com/1lann/smarter-hospital/views/imports"

	"github.com/1lann/smarter-hospital/modules/climate"
	_ "github.com/1lann/smarter-hospital/modules/heartrate"
	"github.com/1lann/smarter-hospital/modules/lights"
	"github.com/1lann/smarter-hospital/modules/proximity"
	"github.com/1lann/smarter-hospital/modules/thermometer"
	"github.com/1lann/smarter-hospital/modules/ultrasonic"
)

var server *core.Server

var webPath = os.Getenv("GOPATH") + "/src/github.com/1lann/smarter-hospital/server"

var connectedModules = make(map[string]bool)
var connectedMutex = new(sync.Mutex)

func moduleSetup() {
	core.SetupModule("lights", "lights1", lights.Settings{
		Pin:               12,
		AnimationDuration: time.Second,
	})

	core.SetupModule("ultrasonic", "ultrasonic1", ultrasonic.Settings{
		TriggerPin: 14,
		EchoPin:    15,
	})

	core.SetupModule("heartrate", "heartrate1")
	core.SetupModule("proximity", "proximity1", proximity.Settings{
		PersonHeight: 100,
	})

	core.SetupModule("thermometer", "thermometer1", thermometer.Settings{
		DeviceID: "28-000004dddaa1",
	})

	core.SetupModule("climate", "climate1", climate.Settings{
		CoolingPin: 11,
		HeatingPin: 13,
		MaxHeating: 255,
		MaxCooling: 255,
	})
}

func main() {
	err := store.Connect(store.ConnectOpts{
		Address:  "127.0.0.1:27017",
		Database: "smarter-hospital",
	})
	if err != nil {
		panic(err)
	}

	moduleSetup()

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

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	notifyServer := notify.NewServer(wsServer)
	logic.Register(r, wsServer, notifyServer)

	server, err = core.NewServer("0.0.0.0:5000")
	if err != nil {
		panic(err)
	}

	r.LoadHTMLFiles(webPath + "/view.tmpl")

	allPages := views.AllPages()
	for pagePath, page := range allPages {
		r.GET(pagePath, pageHandler(page.Title()))
	}

	r.GET("/ws", func(c *gin.Context) {
		wsServer.Handle(c.Request, c.Writer)
	})

	r.GET("/module/info/:moduleid", handleInfo)
	r.POST("/module/action/:moduleid", handleAction)
	r.GET("/module/connected", func(c *gin.Context) {
		var results []string
		connectedMutex.Lock()
		for module := range connectedModules {
			results = append(results, module)
		}
		connectedMutex.Unlock()
		c.JSON(http.StatusOK, results)
	})

	r.GET("/notify/all", func(c *gin.Context) {
		n, err := notifyServer.Notifications()
		if err != nil {
			log.Println("server: notify: all:", err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, n)
		return
	})

	r.GET("/notify/dismiss/:id", func(c *gin.Context) {
		err := notifyServer.Dismiss(c.Param("id"))
		if err != nil {
			log.Println("server: notify: dismiss:", err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.String(http.StatusOK, "")
	})

	r.GET("/notify/call", func(c *gin.Context) {
		err := notifyServer.Push(notify.Notification{
			Alert:    true,
			Heading:  "Ash Ketchum made a nurse call",
			Location: "Ash Ketchum - Room 025",
			Icon:     "red doctor",
			Link:     "/nurse/room",
		})
		if err != nil {
			log.Println(err)
		}

		c.String(http.StatusOK, "")
	})

	r.Static("/static", webPath+"/vendor")

	log.Println("Server is running!")
	r.Run()
}

func pageHandler(title string) func(c *gin.Context) {
	return func(c *gin.Context) {
		if title == "Nurse Controls" {
			c.HTML(http.StatusOK, "view.tmpl", gin.H{
				"Title": title,
				"User": gin.H{
					"firstName": "Nurse",
					"lastName":  "Joy",
				},
			})
		} else {
			c.HTML(http.StatusOK, "view.tmpl", gin.H{
				"Title": title,
				"User": gin.H{
					"firstName": "Ash",
					"lastName":  "Ketchum",
				},
			})
		}

	}
}
