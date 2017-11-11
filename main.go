// Package main provides a game table manager
package main

import (
	"encoding/binary"
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
	"A": 0x77,
	"B": 0x1F,
	"C": 0x4E,
	"D": 0x3D,
	"E": 0x4F,
	"F": 0x47,
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

var nibble = map[byte]byte{
	0x0: 0x7E,
	0x1: 0x30,
	0x2: 0x6D,
	0x3: 0x79,
	0x4: 0x33,
	0x5: 0x5B,
	0x6: 0x5F,
	0x7: 0x70,
	0x8: 0x7F,
	0x9: 0x7B,
	0xA: 0x77,
	0xB: 0x1F,
	0xC: 0x4E,
	0xD: 0x3D,
	0xE: 0x4F,
	0xF: 0x47,
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
	count uint16
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
	countBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(countBytes, c.count)

	c.buf[0] = nibble[(countBytes[0]&0xF0)>>4]
	c.buf[2] = nibble[(countBytes[0]&0x0F)>>0]
	c.buf[4] = nibble[(countBytes[1]&0xF0)>>4]
	c.buf[6] = nibble[(countBytes[1]&0x0F)>>0]

	c.count++

	c.led.SetLEDs(c.buf)
	c.led.Show()
	return
}

func rotate(b byte) byte {
	return b<<1 | b>>7
}
