package src

import "time"

type systemClock struct{}

type Clock interface {
	// return current local time
	Now() time.Time

	// Start "clock"
	NewTicker(time.Duration) *time.Ticker
}

var SystemClock = systemClock{}

func (systemClock) Now() time.Time {
	return time.Now()
}

func (systemClock) NewTicker(duration time.Duration) *time.Ticker {
	return time.NewTicker(duration)
}
