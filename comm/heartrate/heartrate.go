package heartrate

import "encoding/gob"

func init() {
	gob.RegisterName("hre", Event{})
}

// Event represents the event value for a heart rate monitor.
type Event struct {
	BPM     float64
	Contact bool
}
