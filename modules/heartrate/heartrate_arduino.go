// +build !js

package heartrate

import (
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

	peaks := make([]time.Duration, 5)
	var lastPeak time.Time

	isAbove := false
	newData := false

	for range time.Tick(time.Millisecond * 10) {
		val, err := arduino.Adaptor.AnalogRead(pin)
		if err != nil {
			log.Println("heartrate read:", err)
			continue
		}

		if !isAbove && (val > m.Settings.PeakThreshold) {
			isAbove = true
		} else if isAbove && (val <= m.Settings.PeakThreshold) {
			newData = true
			peaks = append(peaks[1:], time.Since(lastPeak))
			lastPeak = time.Now()
			isAbove = false
		}

		if time.Since(lastPeak) > (time.Second * 2) {
			// log.Println("heartrate: no contact?")
		}

		if newData {
			newData = false
			valid := true
			var sum time.Duration
			for _, peak := range peaks {
				if peak == 0 {
					// log.Println("heartrate: not enough data")
					valid = false
					break
				}

				if peak > (time.Second * 2) {
					// log.Println("heartrate: no contact peak?")
					valid = false
					break
				}

				sum += peak
			}

			if valid {
				log.Println("BPM:", 60.0/(sum.Seconds()/float64(len(peaks))))
			}
		}
	}
}
