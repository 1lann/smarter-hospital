package drivers

import (
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/convertors/mcp3008"
	"gobot.io/x/gobot/platforms/raspi"
)

// Analog represents the analog driver to MCP3008 ADC over SPI.
var Analog *mcp3008.MCP3008

// I2CBus represents the I2C bus.
var I2CBus embd.I2CBus

var GoBot *raspi.Adaptor

// Connect starts the analog MCP3008 driver and the I2C driver.
func Connect() error {
	spiBus := embd.NewSPIBus(embd.SPIMode0, 0, 3600000, 8, 0)
	Analog = mcp3008.New(mcp3008.SingleMode, spiBus)
	I2CBus = embd.NewI2CBus(1)
	GoBot = raspi.NewAdaptor()

	return nil
}
