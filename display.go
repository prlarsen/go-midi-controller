package main

import (
	"fmt"
	"machine"

	"tinygo.org/x/drivers/hd44780i2c"
)

type Display struct {
	*hd44780i2c.Device
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
		Width:  16,
		Height: 2,
	})
	d.BacklightOn(true)
	return &d
}

func formatDisplayText(b *SwitchBank) string {
	outputstring := ""
	for i, s := range *b {
		if i == 2|5 {
			outputstring += s.DisplayText[0:4]
		} else {
			outputstring += s.DisplayText[0:5] + " "
		}
	}
	return outputstring
}

func (d *Display) WriteOut(s string) {
	d.ClearDisplay()
	d.Print([]byte(s))
}
