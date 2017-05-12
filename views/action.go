package views

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/1lann/smarter-hospital/core"
	"github.com/gopherjs/gopherjs/js"
)

// Some of the possible error values from Do()
var (
	ErrNotFound   = errors.New("views: action not found")
	ErrBadRequest = errors.New("views: action bad request")
)

// Do performs an action to the server over HTTP.
func Do(moduleID string, val interface{}) (string, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	var scheme string
	if js.Global.Get("location").Get("protocol").String() == "https:" {
		scheme = "https://"
	} else {
		scheme = "http://"
	}

	resp, err := http.Post(scheme+
		js.Global.Get("location").Get("host").String()+"/action/"+moduleID,
		"application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return "", ErrBadRequest
	case http.StatusNotFound:
		return "", ErrNotFound
	case http.StatusInternalServerError:
		break
	case http.StatusOK:
		break
	default:
		return "", errors.New("Unknown response")
	}

	var response core.Result

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&response)
	if err != nil {
		return "", err
	}

	if response.Successful {
		return response.Message, nil
	}

	return "", errors.New(response.Message)
}
