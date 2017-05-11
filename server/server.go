package main

import (
	"log"
	"net/http"
	"os"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/logic"
	"github.com/1lann/smarter-hospital/store"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gin-gonic/gin"

	"github.com/1lann/smarter-hospital/views"
	_ "github.com/1lann/smarter-hospital/views/imports"

	_ "github.com/1lann/smarter-hospital/modules/ping"
)

var server *core.Server

var webPath = os.Getenv("GOPATH") + "/src/github.com/1lann/smarter-hospital/server"

func main() {
	err := store.Connect(store.ConnectOpts{
		Address:  "127.0.0.1:27017",
		Database: "smarter-hospital",
	})
	if err != nil {
		panic(err)
	}

	core.SetupModule("ping", "ping1")

	wsServer := ws.NewServer()
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

	r.POST("/action/:moduleid", handleAction)

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
