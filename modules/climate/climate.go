package climate

// State ...
type State int

const (
	StateOff = iota
	StateCooling
	StateHeating
)

// Settings ...
type Settings struct {
	CoolingPin int
	HeatingPin int

	MaxCooling byte
	MaxHeating byte
}

// Action ...
type Action struct {
	State     State
	Intensity float64
}

// Event ...
type Event struct {
	State     State
	Intensity float64
}
