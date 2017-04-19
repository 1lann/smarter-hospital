package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/1lann/smarter-hospital/core"
	"github.com/gin-gonic/gin"
)

func handleAction(c *gin.Context) {
	var actionPost struct {
		ModuleID string
		Data     json.RawMessage
	}

	c.BindJSON(&actionPost)

	actionData, err := core.ActionValue(actionPost.ModuleID)
	if err != nil {
		log.Println("handleWS error while retrieving", actionPost.ModuleID,
			":", err)
	}

	err = json.Unmarshal(actionPost.Data, &actionData)
	if err != nil {
		log.Println("handleWS error while unmarshalling data:", err)
	}

	resp, err := server.Do(actionPost.ModuleID, actionData)
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
