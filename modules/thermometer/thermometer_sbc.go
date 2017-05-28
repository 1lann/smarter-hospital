// +build !js

package thermometer

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/1lann/smarter-hospital/core"
)

// Module ...
type Module struct {
	ID string
	Settings

	LastEvent Event
}

func init() {
	core.RegisterModule(Module{})
}

// HandleEvent ...
func (m *Module) HandleEvent(evt Event) {
	m.LastEvent = evt
}

// Info ...
func (m *Module) Info() Event {
	return m.LastEvent
}

type tempRatio struct {
	temp  float64
	ratio float64
}

var tempRatios []tempRatio

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	t := time.NewTicker(time.Second)
	for range t.C {
		file, err := os.Open("/sys/bus/w1/devices/" + m.Settings.DeviceID + "/w1_slave")
		if err != nil {
			client.Error(m.ID, err)
			continue
		}

		data, err := ioutil.ReadAll(file)
		file.Close()
		if err != nil {
			client.Error(m.ID, err)
			continue
		}

		lines := strings.Split(string(data), "\n")
		message := strings.Split(lines[0], " ")
		if message[len(message)-1] != "YES" {
			client.Error(m.ID, errors.New("message does not end with YES"))
			continue
		}

		dataLine := strings.Split(lines[1], " ")
		if len(dataLine) < 10 {
			client.Error(m.ID, errors.New("message not long enough?"))
			continue
		}

		temp, err := strconv.Atoi(strings.TrimPrefix(dataLine[9], "t="))
		if err != nil {
			client.Error(m.ID, err)
			continue
		}

		log.Println("temp:", float64(temp)/1000)

		client.Emit(m.ID, Event{
			Temperature: float64(temp) / 1000,
		})
	}
}
