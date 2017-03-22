package proto

// Constant event message values.
const (
	HealthCheckMsg = "health_check"
	EventMsg       = "event"
	ActionMsg      = "action"
	AuthMsg        = "auth"
)

// Event represents an event from a sensor to the system.
type Event struct {
	Name  string
	Value interface{}
}

// Action represents an action issued by the system to an actuator.
type Action struct {
	Name  string
	Value interface{}
}

// Result represents the result of an action or the emission of an event.
type Result struct {
	Successful bool
	Message    string
}

// AuthRequest represents a request to authenticate with the system.
type AuthRequest struct {
	ID string
}

// AuthResponse represents an authentication response from the system.
type AuthResponse struct {
	Authenticated bool
}

// ActionHandler represents a rpc2 handler for actions.
type ActionHandler func(action Action) error
