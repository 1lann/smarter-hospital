// +build !js

package thermistor

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/1lann/smarter-hospital/arduino"
	"github.com/1lann/smarter-hospital/core"
)

type tempRatio struct {
	temp  float64
	ratio float64
}

var tempRatios []tempRatio

// PollEvents ...
func (m *Module) PollEvents(client *core.Client) {
	var f float64
	for f = 273; f <= 320; f += 0.1 {
		tempRatios = append(tempRatios,
			tempRatio{
				temp:  f,
				ratio: getR25Ratio(f),
			})
	}

	pin := strconv.Itoa(m.Settings.Pin + 4)

	t := time.NewTicker(time.Second)
	for range t.C {
		val, err := arduino.Adaptor.AnalogRead(pin)
		if err != nil {
			log.Println("thermistor: read:", err)
			continue
		}

		temp := m.getTemperature(m.getResistance((float64(val) / float64(1023)) * 5.0))

		log.Println("temp:", temp-273.15)

		client.Emit(m.ID, Event{
			Temperature: temp,
		})
	}
}

func (m *Module) getResistance(v float64) float64 {
	i := (m.Settings.SupplyVoltage - v) / m.Settings.FixedResistor
	return (m.Settings.SupplyVoltage - i*m.Settings.FixedResistor) / i
}

func (m *Module) getTemperature(r float64) float64 {
	for _, temp := range tempRatios {
		if temp.ratio*m.Settings.R25 <= r {
			return temp.temp
		}
	}

	return 0
}

func getR25Ratio(temp float64) float64 {
	return math.Exp((-1.4122478e1) + (4.4136033e3)/temp + (-2.9034189e4)/math.Pow(temp, 2) - (9.3875035e6)/math.Pow(temp, 3))
}
