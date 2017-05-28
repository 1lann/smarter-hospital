package heartrate

// Settings ...
type Settings struct{}

// Action ...
type Action struct{}

// Event ...
type Event struct {
	Contact     bool
	Calculating bool
	BPM         float64
}
