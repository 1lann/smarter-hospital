package lights

import "time"

// Settings ...
type Settings struct {
	Pin               int
	AnimationDuration time.Duration
}

// Action ...
type Action struct {
	State int // 0 = off, 255 = full
}

// Event ...
type Event struct {
	NewState int // 0 = off, 255 = full
	// Time     time.Time
}
