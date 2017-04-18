package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/1lann/smarter-hospital/comm"
	"github.com/gin-gonic/gin"
)

func handleAction(c *gin.Context) {
	var actionPost struct {
		Action string
		Device string
		Data   json.RawMessage
	}

	c.BindJSON(&actionPost)

	actionData, err := core.ActionValue(actionPost.Action)
	if err != nil {
		log.Println("handleWS error while retrieving", actionPost.Action,
			":", err)
	}

	err = json.Unmarshal(actionPost.Data, &actionData)
	if err != nil {
		log.Println("handleWS error while unmarshalling data:", err)
	}

	resp, err := server.Do(actionPost.Device, actionPost.Action, actionData)
	log.Println("server: action post:", err)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Success": false,
			"Message": resp,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Success": true,
		"Message": resp,
	})
}
