package heartrate

import (
	"log"
	"time"

	"github.com/1lann/smarter-hospital/store"
	"gopkg.in/mgo.v2/bson"
)

const timeframe = 10 * time.Minute

func init() {
	core.RegisterModule(Module{})
}

// Module represents the module.
type Module struct {
	ID string

	Event
	Settings

	LastBPM Event
}

// Event represents the event value for a heart rate monitor.
type Event struct {
	BPM     float64
	Contact bool
}

// Settings represents the settings of the module.
type Settings struct {
	Pin int
}

// Info retrieves the current and historical information of module.
func (c *Context) Info() interface{} {
	var history []Event

	err := store.C("hr_history").Find(bson.M{
		"Time": bson.M{
			"$gt": time.Now().Add(-timeframe),
		},
	}).Sort("Time").All(&history)
	if err != nil {
		return err
	}

	return struct {
		Now     Event
		History []Event
	}{
		c.LastBPM,
		history,
	}
}

// HandleEvent handles an event on the server.
func (c *Context) HandleEvent(evt Event) error {
	c.LastBPM = evt

	store.C("hr_history").EnsureIndexKey("Time")
	err := store.C("hr_history").Insert(bson.M{
		"BPM":  evt.BPM,
		"Time": time.Now(),
	})

	if err != nil {
		log.Println("heartrate: store:", err)
	}

	return nil
}
