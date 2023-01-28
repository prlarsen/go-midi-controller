package main

import (
	"fmt"
	"machine"
	"strings"

	"tinygo.org/x/drivers/hd44780i2c"
)

const DISPWIDTH uint8 = 16
const DISPHEIGHT uint8 = 2

type Display struct {
	*hd44780i2c.Device
}

func (d *Display) WriteOut(s []byte) {
	d.ClearDisplay()
	d.Print(s)
}

func initDisplay() *hd44780i2c.Device {
	i2c := machine.I2C1
	err := i2c.Configure(machine.I2CConfig{
		SCL: machine.GP11,
		SDA: machine.GP10,
	})
	if err != nil {
		fmt.Println(err)
	}
	d := hd44780i2c.New(i2c, 0)
	d.Configure(hd44780i2c.Config{
		Width:  DISPWIDTH,
		Height: DISPHEIGHT,
	})
	d.BacklightOn(true)
	return &d
}

func trimText(s string, max int) string {
	if len(s) >= max {
		return s[:max]
	} else {
		return s
	}
}

// TODO: make this function more generic so that it will
// compute the proper length and spacing of text based
// on number of inputs and screen dimensions
func formatDisplayText(b SwitchBank, banknum uint8) string {
	outputstring := ""
	for i := 3; i <= 5; i++ {
		if i == 5 {
			outputstring += fmt.Sprintf("BNK%d\n", banknum+1)
		} else {
			t := trimText(b[i].DisplayText, 4)
			padding := strings.Repeat(" ", 4-len(t))
			outputstring += t + padding + "  "
		}
	}
	for i := 0; i < 2; i++ {
		t := trimText(b[i].DisplayText, 4)
		padding := strings.Repeat(" ", 4-len(t))
		outputstring += t + padding + "  "
	}
	t := trimText(b[2].DisplayText, 4)
	outputstring += t
	return outputstring
}
