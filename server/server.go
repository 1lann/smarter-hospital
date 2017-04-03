package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/1lann/multitemplate"
	"github.com/1lann/smarter-hospital/comm"
	"github.com/1lann/smarter-hospital/comm/blinker"
	"github.com/gin-gonic/gin"
)

var server *comm.Server

var webPath = os.Getenv("GOPATH") + "/src/github.com/1lann/smarter-hospital/server"

func main() {
	var err error
	server, err = comm.NewServer("0.0.0.0:5000", authHandler, handlers)
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

	r.GET("/hello/:msg", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, %v", c.Param("msg"))
	})

	r.GET("/blink/:n", func(c *gin.Context) {
		blinkRate, err := strconv.Atoi(c.Param("n"))
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"error": "blink rate must be a number",
			})
			return
		}

		resp, err := server.Do("arduino", "blink", blinker.Action{
			Rate: blinkRate,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not perform action on device: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result":   "OK",
			"response": resp,
		})
	})

	r.Static("/static", webPath+"/static")

	r.HTMLRender = t

	log.Println("Server is running!")
	r.Run()
}
