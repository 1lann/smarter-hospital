package ultrasonic

// Settings ...
type Settings struct {
	TriggerPin       int
	EchoPin          int
	ContactThreshold float64
}

// Action ...
type Action struct{}

// Event ...
type Event struct {
	Contact bool
}
