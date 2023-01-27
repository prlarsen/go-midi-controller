package main

import (
	// "encoding/binary"

	"time"

	// Look at the tinyjson package for creating custom json marshaler

	"machine"
	// "tinygo.org/x/drivers/at24cx"
)

// DEBOUNCE_TIME controls how long to ignore subsequent input
const DEBOUNCE_TIME = 500 * time.Millisecond

// Define the amount of time to determine a held down switch
const SWITCH_HOLD_TIME = 500 * time.Millisecond

// Default sleep duration
const SLEEP_5MS = 5 * time.Millisecond

const OB_LED = machine.LED

// Common value assingment MIDI ON
const MIDI_VALUE_ON uint8 = 0x7f
const MIDI_VALUE_OFF uint8 = 0x00

func switchIsHeld(pin machine.Pin) bool {
	trigger := time.Now()

	for time.Now().Sub(trigger) < SWITCH_HOLD_TIME {
		if !pin.Get() {
			return false
		}
		time.Sleep(SLEEP_5MS)
	}
	return true
}

func selectBank(numOfBanks uint8, currentBank uint8, foward bool) uint8 {
	if foward {
		if currentBank == numOfBanks-1 {
			return 0
		} else {
			return currentBank + 1
		}
	} else {
		if currentBank == 0 {
			return numOfBanks - 1
		} else {
			return currentBank - 1
		}
	}
}

func main() {
	OB_LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	disp := Display{initDisplay()}
	// eeprom := initEEPROM()

	// Initialize controller inputs to first bank
	inputBank := uint8(0)

	ConfigureSwitchPins()

	ctrlInputs := InitDefaultSwitchBanks()
	numOfBanks := uint8(len(ctrlInputs))
	bankSelect := SWITCHPINS[5]

	disp.WriteOut(formatDisplayText(&ctrlInputs[inputBank]))

	for {

		time.Sleep(SLEEP_5MS)

		if bankSelect.Get() {
			inputBank = selectBank(numOfBanks, inputBank, switchIsHeld(bankSelect))
		}

		for _, inp := range ctrlInputs[inputBank] {

			if inp.Pin.Get() {
				inp.sendControl()
			}
		}
	}
}
