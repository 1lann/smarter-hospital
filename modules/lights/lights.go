package lights

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/store"
)

func init() {
	core.RegisterModule(Module{})
}

type Module struct {
	ID string
}

type Action struct{}

type Event struct{}

func (m *Module) HandleAction(client *core.Client, act Action) error {
	log.Println("I received an action")
	return nil
}

func (m *Module) PollEvents(client *core.Client) {
	for {
		for i := 0; i < 128; i++ {
			// if i > 64 {
			// 	arduino.Adaptor.PwmWrite("13", byte(128-i))
			// } else {
			// 	arduino.Adaptor.PwmWrite("13", byte(i))
			// }

			time.Sleep(time.Millisecond * 20)
		}

		log.Println("emitting an event")
		client.Emit(m.ID, Event{})
	}
}

func (m *Module) HandleEvent(evt Event) {
	log.Println("I received an event")
	store.C("lights_test").Insert(bson.M{
		"word": "memes",
		"time": time.Now(),
	})
}

func (m *Module) Info() Event {
	return Event{}
}
