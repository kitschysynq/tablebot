// Package main provides a game table manager
package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"

	"github.com/kitschysynq/tablebot/ht16k33"
)

var (
	AllOn = []byte{
		0xFF, 0xFF,
		0xFF, 0xFF,
		0xFF, 0xFF,
		0xFF, 0xFF,
		0xFF, 0xFF,
		0xFF, 0xFF,
		0xFF, 0xFF,
		0xFF, 0xFF,
	}
	AllOff = []byte{
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
		0x00, 0x00,
	}
)

func main() {
	r := raspi.NewAdaptor()
	ctl := &controller{
		buf: []byte{
			0x01, 0x01,
			0x01, 0x01,
			0x01, 0x01,
			0x01, 0x01,
			0x01, 0x01,
			0x01, 0x01,
			0x01, 0x01,
			0x01, 0x01,
		},
		led: ht16k33.NewHT16K33Driver(r),
	}

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			fmt.Printf("state: %#0x\n", ctl.buf[0])
			ctl.Toggle()
		})
	}

	robot := gobot.NewRobot("tablebot",
		[]gobot.Connection{r},
		[]gobot.Device{ctl},
		work,
	)

	robot.Start()
}

type controller struct {
	on  bool
	buf []byte
	led *ht16k33.HT16K33Driver
}

func (c *controller) Name() string                 { return c.led.Name() }
func (c *controller) SetName(s string)             { c.led.SetName(s) }
func (c *controller) Connection() gobot.Connection { return c.led.Connection() }
func (c *controller) Halt() error                  { return c.led.Halt() }

// Start initiates the Driver
func (c *controller) Start() error {
	err := c.led.Start()
	if err != nil {
		return err
	}

	c.led.SetLEDs(c.buf)
	c.led.Show()
	return nil
}

func (c *controller) Toggle() {
	/*
		if c.on {
			c.led.SetLEDs(AllOff)
			c.on = false
		} else {
			c.led.SetLEDs(AllOn)
			c.on = true
		}
	*/
	for i := range c.buf {
		c.buf[i] = rotate(c.buf[i])
		if c.buf[i] > 63 {
			c.buf[i] = 1
		}
	}

	c.led.SetLEDs(c.buf)
	c.led.Show()
	return
}

func rotate(b byte) byte {
	return b<<1 | b>>7
}
