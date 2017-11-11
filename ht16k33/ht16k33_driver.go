package ht16k33

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

const (
	ht16k33DefaultAddress = 0x70
)

type HT16K33Int byte

const (
	HT16K33InterruptActiveLow HT16K33Int = iota
	HT16K33InterruptActiveHigh
)

const (
	ht16k33DisplayData  = 0x00
	ht16k33SystemSetup  = 0x20
	ht16k33KeyData      = 0x40
	ht16k33IntFlag      = 0x60
	ht16k33DisplaySetup = 0x80
	ht16k33RowInt       = 0xA0
	ht16k33Dimming      = 0xE0
)

type HT16K33Mode byte

const (
	HT16K33Row           HT16K33Mode = 0x00
	HT16K33IntActiveLow              = 0x01
	HT16K33IntActiveHigh             = 0x03
)

// HT16K33Driver is the gobot driver for the Holtek HT16K33 LED Controller Driver
//
// Device datasheet: http://www.ti.com/lit/ds/symlink/drv2605l.pdf
//
// Basic use:
//
//  led_key := i2c.NewHT16K33Driver(adaptor)
//  led_key.SetLEDs([]byte{1, 13})
//  led_key.Show()

var _ gobot.Device = (*HT16K33Driver)(nil)

//
type HT16K33Driver struct {
	name       string
	connector  i2c.Connector
	connection i2c.Connection
	i2c.Config

	leds []byte
	keys []byte
}

// NewHT16K33Driver creates a new driver for the HT16K33 device.
//
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewHT16K33Driver(conn i2c.Connector, options ...func(i2c.Config)) *HT16K33Driver {
	driver := &HT16K33Driver{
		name:      gobot.DefaultName("HT16K33"),
		connector: conn,
		Config:    i2c.NewConfig(),
		leds:      make([]byte, 16),
		keys:      make([]byte, 6),
	}

	for _, option := range options {
		option(driver)
	}

	return driver
}

// Name returns the name of the device.
func (d *HT16K33Driver) Name() string {
	return d.name
}

// SetName sets the name of the device.
func (d *HT16K33Driver) SetName(name string) {
	d.name = name
}

// Connection returns the connection of the device.
func (d *HT16K33Driver) Connection() gobot.Connection {
	return d.connector.(gobot.Connection)
}

// Start initializes the device.
func (d *HT16K33Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetDefaultBus())
	address := d.GetAddressOrDefault(ht16k33DefaultAddress)

	d.connection, err = d.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	// Start the system clock
	err = d.connection.WriteByte(ht16k33SystemSetup | 0x01)
	if err != nil {
		return err
	}

	// Turn the display on
	err = d.connection.WriteByte(ht16k33DisplaySetup | 0x01)
	if err != nil {
		return err
	}

	return nil
}

// Halt terminates the Driver
func (d *HT16K33Driver) Halt() error {
	return nil
}

// SetRowInt sets the device to either use all available rows for keyscan,
// or to trade one row of keyscan for an input interrupt pin
func (d *HT16K33Driver) SetRowInt(mode HT16K33Mode) (err error) {
	return d.connection.WriteByte(byte(ht16k33RowInt | mode))
}

// SetDisplay turns the display on or off and sets blink frequency
//func (d *HT16K33Driver) SetDisplay(f HT16K33DisplaySetupFlag) error {
//	return d.connection.WriteByte(ht16k33DisplaySetup | f)
//}

// SetLEDs overwrites the state of the LEDs
func (d *HT16K33Driver) SetLEDs(leds []byte) error {
	if len(leds) != 16 {
		return fmt.Errorf("need 16 bytes for led data")
	}
	copy(d.leds, leds)
	return nil
}

// Show sends the current display buffer to the device output
func (d *HT16K33Driver) Show() error {
	_, err := d.connection.Write(append([]byte{ht16k33DisplayData}, d.leds...))
	return err
}

// Dim sets the brightness of the display output via built-in PWM
func (d *HT16K33Driver) Dim(dutycycle uint8) error {
	if dutycycle > 0xF {
		return fmt.Errorf("duty cycle out of range")
	}
	return d.connection.WriteByte(ht16k33Dimming | dutycycle)
}

func (d *HT16K33Driver) ReadKeyData() []uint8 {
	return nil
}
