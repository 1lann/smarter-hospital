package views

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/1lann/smarter-hospital/core"
)

// Address stores the address of the website including the scheme, based on the
// location provided by the browser.
var Address string

// Some of the possible error values from Do()
var (
	ErrNotFound      = errors.New("views: action not found")
	ErrBadRequest    = errors.New("views: action bad request")
	ErrInternalError = errors.New("views: internal error")
)

// ModuleDo performs a device action to the server over HTTP.
func ModuleDo(moduleID string, val interface{}) (string, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(Address+"/module/action/"+moduleID,
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

// ModuleInfo requests information regarding a module, the response argument
// should be a pointer where the response will be decoded to.
func ModuleInfo(moduleID string, response interface{}) error {
	resp, err := http.Get(Address + "/module/info/" + moduleID)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusInternalServerError:
		return ErrInternalError
	case http.StatusOK:
		break
	default:
		return errors.New("Unknown response")
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(response)
	if err != nil {
		return err
	}

	return nil
}

// ConnectedModules requests the server over HTTP for the currently connected
// modules.
func ConnectedModules() ([]string, error) {
	resp, err := http.Get(Address + "/module/connected")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var response []string
	err = dec.Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
