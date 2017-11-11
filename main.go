// Package main provides a game table manager
package main

import (
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

var charset = map[string]byte{
	"0": 0x7E,
	"1": 0x30,
	"2": 0x6D,
	"3": 0x79,
	"4": 0x33,
	"5": 0x5B,
	"6": 0x5F,
	"7": 0x70,
	"8": 0x7F,
	"9": 0x7B,
}

var numerals = map[int]byte{
	0: 0x7E,
	1: 0x30,
	2: 0x6D,
	3: 0x79,
	4: 0x33,
	5: 0x5B,
	6: 0x5F,
	7: 0x70,
	8: 0x7F,
	9: 0x7B,
}

func main() {
	r := raspi.NewAdaptor()
	ctl := &controller{
		buf: []byte{
			0x02, 0x02,
			0x02, 0x02,
			0x02, 0x02,
			0x02, 0x02,
			0x02, 0x02,
			0x02, 0x02,
			0x02, 0x02,
			0x02, 0x02,
		},
		led: ht16k33.NewHT16K33Driver(r),
	}

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
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
	count int
	buf   []byte
	led   *ht16k33.HT16K33Driver
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
	for i := range c.buf {
		c.buf[i] = numerals[c.count]
		c.count++
		c.count %= 10
	}

	c.led.SetLEDs(c.buf)
	c.led.Show()
	return
}

func rotate(b byte) byte {
	return b<<1 | b>>7
}
