package main

import (
	"time"

	"machine"
)

// DEBOUNCE_TIME controls how long to ignore subsequent input
const DEBOUNCE_TIME = 500 * time.Millisecond

// Define the amount of time to determine a held down switch
const SWITCH_HOLD_TIME = 500 * time.Millisecond

// Default sleep duration
const SLEEP_5MS = 5 * time.Millisecond

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

func selectBank(numOfBanks uint8, currentBank uint8, reverse bool) uint8 {
	if reverse {
		if currentBank == 0 {
			return numOfBanks - 1
		} else {
			return currentBank - 1
		}
	} else {
		if currentBank == numOfBanks-1 {
			return 0
		} else {
			return currentBank + 1
		}
	}
}

func main() {
	disp := Display{initDisplay()}
	// eeprom := initEEPROM()

	inputBank := uint8(0) // Initialize controller inputs to first bank

	ConfigureSwitchPins()

	ctrlInputs := InitDefaultSwitchBanks()
	numOfBanks := uint8(len(ctrlInputs))
	bankSelect := SWITCHPINS[5]

	disp.WriteOut([]byte(formatDisplayText(ctrlInputs[inputBank], inputBank)))

	for {
		time.Sleep(SLEEP_5MS)

		if bankSelect.Get() {
			inputBank = selectBank(numOfBanks, inputBank, switchIsHeld(bankSelect))
			disp.WriteOut([]byte(formatDisplayText(ctrlInputs[inputBank], inputBank)))
			time.Sleep(DEBOUNCE_TIME)
		}

		for _, inp := range ctrlInputs[inputBank] {

			if inp.Pin.Get() {
				inp.sendControl()
				time.Sleep(DEBOUNCE_TIME)
			}
		}
	}
}
