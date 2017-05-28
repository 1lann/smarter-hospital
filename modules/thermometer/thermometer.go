package thermometer

// Settings ...
type Settings struct {
	DeviceID string
}

// Action ...
type Action struct{}

// Event ...
type Event struct {
	Temperature float64 // In Celsius
}
