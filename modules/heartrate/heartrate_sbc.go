// +build !js

package heartrate

import (
	"math"
	"sort"
	"time"

	"github.com/1lann/max30105"
	"github.com/1lann/smarter-hospital/core"
	"github.com/1lann/smarter-hospital/pi/drivers"
	"github.com/mjibson/go-dsp/spectral"
)

// 25 samples per second

// Module ...
type Module struct {
	ID string
	Settings

	LastBPM [5]float64
}

func init() {
	core.RegisterModule(Module{})
}

// HandleEvent ...
func (m *Module) HandleEvent(evt Event) {}

type pwelchPoint struct {
	Frequency float64
	Amplitude float64
}

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	d := max30105.NewDriver(drivers.I2CBus)
	err := d.Setup()
	if err != nil {
		client.Error(m.ID, err)
		return
	}

	thisPeriod := make([]float64, 256)
	light := false
	contact := false

	lastCalculation := time.Now()
	ticker := time.NewTicker(time.Millisecond * 100)
	for range ticker.C {
		// log.Println(d.ReadTemperature())
		samples, err := d.ReadSamples()
		if err != nil {
			client.Error(m.ID, err)
			return
		}

		data := make([]float64, len(samples))
		for i, sample := range samples {
			// fmt.Println(sample.Red)
			data[i] = float64(sample.Red)
		}

		// fmt.Println("Samples:", len(samples))

		if len(samples) > 0 {
			if samples[0].IR > 20000 && !light {
				contact = true
				thisPeriod = make([]float64, 256)
				m.LastBPM = [5]float64{0, 0, 0, 0, 0}

				client.Emit(m.ID, Event{
					Contact:     true,
					Calculating: true,
					BPM:         0,
				})

				if d.SetRedAmplitude(0x1F) == nil && d.SetGreenAmplitude(0x1F) == nil {
					light = true
				}
			}

			if samples[0].IR <= 20000 && samples[0].Green <= 20000 && light {
				contact = false

				client.Emit(m.ID, Event{
					Contact:     false,
					Calculating: false,
					BPM:         0,
				})

				if d.SetRedAmplitude(0) == nil && d.SetGreenAmplitude(0) == nil {
					light = false
				}
			}
		}

		thisPeriod = append(thisPeriod[len(samples):], data...)

		if (time.Since(lastCalculation) > time.Second) && (thisPeriod[0] != 0) && contact {
			var topPoints []pwelchPoint

			for i := 128; i <= 256; i += 2 {
				pxx, freqs := spectral.Pwelch(thisPeriod[len(thisPeriod)-i:],
					25, &spectral.PwelchOptions{
						NFFT: i,
					})

				var points []pwelchPoint

				for i := 0; i < len(pxx); i++ {
					if freqs[i] >= 0.9 && freqs[i] < 5.0 {
						// fmt.Printf("%v\t%v\n", freqs[i], pxx[i])
						points = append(points, pwelchPoint{
							Frequency: freqs[i],
							Amplitude: pxx[i],
						})
					}
				}

				sort.Slice(points, func(i int, j int) bool {
					return points[i].Amplitude > points[j].Amplitude
				})

				top := points[0]

				for i := 1; i <= 5; i++ {
					if (math.Abs((top.Frequency/2)-points[i].Frequency) < 0.2) &&
						((points[i].Amplitude / top.Amplitude) > 0.9) {
						top = points[i]
					}
				}

				topPoints = append(topPoints, top)
			}

			// os.Exit(0)

			modeMap := make(map[int]int)

			for i := 9; i < 50; i++ {
				for _, point := range topPoints {
					if math.Abs(point.Frequency-float64(i)/10) <= 0.15 {
						modeMap[i] = modeMap[i] + 1
					}
				}
			}

			modeFrequency := 0
			numMode := 0
			for freq, n := range modeMap {
				if n > numMode {
					numMode = n
					modeFrequency = freq
				}
			}

			totalPoints := 0.0
			numPoints := 0

			for _, point := range topPoints {
				if math.Abs(point.Frequency-float64(modeFrequency)/10) <= 0.15 {
					totalPoints += point.Frequency
					numPoints++
				}
			}

			bpm := (totalPoints / float64(numPoints)) * 60.0

			m.LastBPM[0] = m.LastBPM[1]
			m.LastBPM[1] = m.LastBPM[2]
			m.LastBPM[2] = m.LastBPM[3]
			m.LastBPM[3] = m.LastBPM[4]
			m.LastBPM[4] = bpm

			if m.LastBPM[0] == 0 {
				client.Emit(m.ID, Event{
					Contact:     true,
					Calculating: false,
					BPM:         bpm,
				})

				lastCalculation = time.Now()
				continue
			}

			client.Emit(m.ID, Event{
				Contact:     true,
				Calculating: false,
				BPM: (m.LastBPM[0] + m.LastBPM[1] + m.LastBPM[2] +
					m.LastBPM[3] + m.LastBPM[4]) / 5,
			})

			lastCalculation = time.Now()
		} else if (time.Since(lastCalculation) > time.Second) &&
			(contact && thisPeriod[0] == 0) {
			client.Emit(m.ID, Event{
				Contact:     true,
				Calculating: true,
				BPM:         0,
			})

			lastCalculation = time.Now()
		}
	}
}
