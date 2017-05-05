package dummy

import (
	"log"

	"github.com/1lann/smarter-hospital/core"
)

func init() {
	core.RegisterModule(Module{})
}

type Module struct {
	ID string
	Settings
	State bool
}

type Settings struct {
	HelloWorld string
}

type Action struct {
	SetState bool
}

type Event struct {
	CurrentState bool
}

func (m *Module) HandleAction(client *core.Client, act Action) error {
	log.Println("I received an action, it says:", act)
	log.Println("OK, responding with event to confirm")
	log.Println("I would like to say HelloWorld is set to:", m.Settings.HelloWorld)

	client.Emit(m.ID, Event{
		CurrentState: act.SetState,
	})

	log.Println("Verify: emitted")

	return nil
}

func (m *Module) PollEvents(client *core.Client) {

}

func (m *Module) HandleEvent(evt Event) {
	log.Println("I received an event, it says:", evt)
	log.Println("I would like to say HelloWorld is set to:",
		m.Settings.HelloWorld)
	m.State = true
}

func (m *Module) Info() Event {
	return Event{
		CurrentState: m.State,
	}
}
