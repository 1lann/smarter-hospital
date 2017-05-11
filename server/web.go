package main

import (
	"log"
	"net/http"
	"reflect"

	"github.com/1lann/smarter-hospital/core"
	"github.com/gin-gonic/gin"
)

func handleAction(c *gin.Context) {
	moduleID := c.Param("moduleid")
	actionType, err := core.ActionType(moduleID)
	if err != nil {
		log.Println("server: error while retrieving action value:",
			moduleID, ":", err)
		c.JSON(http.StatusNotFound, core.Result{
			Successful: false,
			Message:    "No such action",
		})
		return
	}

	actionValue := reflect.New(actionType)

	err = c.BindJSON(actionValue.Interface())
	if err != nil {
		log.Println("server: error while unmarshalling action data:",
			err)
		c.JSON(http.StatusBadRequest, core.Result{
			Successful: false,
			Message:    "Failed to decode action",
		})
	}

	resp, err := server.Do(moduleID, actionValue.Elem().Interface())

	if err != nil {
		log.Println("server: action post:", err)
		c.JSON(http.StatusInternalServerError, core.Result{
			Successful: false,
			Message:    resp,
		})
		return
	}

	c.JSON(http.StatusOK, core.Result{
		Successful: true,
		Message:    resp,
	})
}
