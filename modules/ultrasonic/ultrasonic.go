package ultrasonic

// Settings ...
type Settings struct {
	TriggerPin int
	EchoPin    int
}

// Action ...
type Action struct{}

// Event ...
type Event struct {
	Contact bool
}
