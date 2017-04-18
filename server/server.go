package main

import (
	"log"
	"net/http"
	"os"

	"github.com/1lann/multitemplate"
	"github.com/1lann/smarter-hospital/comm"
	"github.com/1lann/smarter-hospital/ws"
	"github.com/gin-gonic/gin"
)

var server *core.Server

var webPath = os.Getenv("GOPATH") + "/src/github.com/1lann/smarter-hospital/server"

func main() {
	var err error
	server, err = core.NewServer("0.0.0.0:5000", authHandler, handlers)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	t := multitemplate.New()
	t.SetDelimiter("[[", "]]")
	t.AddFromFiles("index", webPath+"/views/index.tmpl")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		ws.Handle(c.Request, c.Writer)
	})

	r.POST("/action/:action", handleAction)

	r.Static("/static", webPath+"/static")

	r.HTMLRender = t

	log.Println("Server is running!")
	r.Run()
}
