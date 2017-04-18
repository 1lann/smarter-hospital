package lights

import (
	"time"

	"github.com/1lann/smarter-hospital/states"
	"gopkg.in/mgo.v2/bson"
)

const timeframe = 10 * time.Minute

func init() {
	states.RegisterStateRetriever(retriever)
	core.RegisterAction(Action{})
}

var lastKnownStates map[string]bool

// Action represents the action value for lighting control.
type Action struct {
	LightID string
	State   bool
}

func retriever() interface{} {
	var history []Event

	err := store.C("light_history").Find(bson.M{
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
		lastBPM,
		history,
	}
}

func store(act Action) {

}
