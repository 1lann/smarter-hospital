// +build !js

package heartrate

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"
)

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	// For some reason pin 4 = analog pin 0
	pin := strconv.Itoa(4 + m.Settings.Pin)
	noContactCount := 0

	peaks := make([]time.Duration, 10)
	var lastPeak time.Time

	isAbove := false
	newData := false

	for range time.Tick(time.Millisecond * 10) {
		val, err := arduino.Adaptor.AnalogRead(pin)
		if err != nil {
			log.Println("heartrate read:", err)
			continue
		}

		fmt.Println(val)

		if !isAbove && (val > m.Settings.PeakThreshold) {
			isAbove = true
		} else if isAbove && (val <= m.Settings.PeakThreshold) {
			newData = true
			peaks = append(peaks[1:], time.Since(lastPeak))
			lastPeak = time.Now()
			isAbove = false
		}

		if time.Since(lastPeak) > (time.Second * 2) {
			noContactCount += 10
		}

		if newData {
			newData = false
			valid := true
			var sum time.Duration
			var numAvailable int
			for i := len(peaks) - 1; i >= 0; i-- {
				numAvailable = len(peaks) - i
				peak := peaks[i]
				if peak == 0 {
					if numAvailable < 5 {
						valid = false
					}
					break
				}

				if peak > (time.Second * 2) {
					noContactCount++

					if numAvailable < 5 {
						valid = false
					}
					break
				}

				sum += peak
			}

			if valid {
				noContactCount = 0
				client.Emit("heartrate1", Event{
					Contact: true,
					BPM:     60.0 / (sum.Seconds() / float64(numAvailable)),
				})
			}

			if noContactCount > 2 {
				client.Emit("heartrate1", Event{
					Contact: false,
				})
			}
		}
	}
}
