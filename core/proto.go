package core

// Constant event message values.
const (
	HealthCheckMsg = "health_check"
	EventMsg       = "event"
	ActionMsg      = "action"
	HandshakeMsg   = "handshake"
)

// Event represents an event from a sensor to the system.
type Event struct {
	ModuleID string
	Value    interface{}
}

// Action represents an action issued by the system to an actuator.
type Action struct {
	ModuleID string
	Value    interface{}
}

// Result represents the result of an action or the emission of an event.
type Result struct {
	Successful bool
	Message    string
}

// HandshakeRequest represents a request to handshake with the system.
type HandshakeRequest struct {
	ModuleIDs []string
}

// HandshakeResponse represents a handshake response from the system.
type HandshakeResponse struct {
	Successful     bool
	ModuleSettings map[string]interface{}
}
