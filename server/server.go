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

	"github.com/1lann/smarter-hospital/modules/heartrate"
	"github.com/1lann/smarter-hospital/modules/lights"
	"github.com/1lann/smarter-hospital/modules/ultrasonic"
)

var server *core.Server

var webPath = os.Getenv("GOPATH") + "/src/github.com/1lann/smarter-hospital/server"

var connectedModules = make(map[string]bool)
var connectedMutex = new(sync.Mutex)

func moduleSetup() {
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
		PeakThreshold: 440,
		Pin:           0,
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

	notifyServer := notify.NewServer(wsServer)
	logic.Register(wsServer, notifyServer)

	server, err = core.NewServer("0.0.0.0:5000")
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
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
