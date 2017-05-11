package ping

import "github.com/1lann/smarter-hospital/core"

// Module ...
type Module struct {
	ID string
	Settings

	LastMessage string
}

// Settings ...
type Settings struct{}

// Action ...
type Action struct {
	Message string
}

// Event ...
type Event struct {
	Message string
}

func init() {
	core.RegisterModule(Module{})
}

// HandleAction ...
func (m *Module) HandleAction(client *core.Client, act Action) error {
	client.Emit(m.ID, Event{
		Message: act.Message,
	})

	return nil
}

// HandleEvent ...
func (m *Module) HandleEvent(evt Event) {
	m.LastMessage = evt.Message
}

// Info ...
func (m *Module) Info() string {
	return m.LastMessage
}
