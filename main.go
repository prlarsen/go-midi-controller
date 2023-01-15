package main

import (
	"time"

	"machine"
	"machine/usb/midi"
)

// DEBOUNCE_TIME controls how long to ignore subsequent input
const DEBOUNCE_TIME = 500 * time.Millisecond

// Define the amount of time to determine a held down switch
const SWITCH_HOLD_TIME = 500 * time.Millisecond

// Default sleep duration
const SLEEP_5MS = 5 * time.Millisecond

// Number of Input Banks
const NUM_INPUT_BANKS = 3

// Setup onboard LED
const OB_LED = machine.LED

// First Control Code defined
const START_CC = uint8(20)

// Common value assingment MIDI ON
const MIDI_VALUE_ON uint8 = 0x7f

type InputType uint8

const (
	SWITCH_ONESHOT InputType = iota
	SWITCH_HOLD
	SWITCH_BANKSELECT // Acts as a SWITCH_HOLD but designated for bank selection
)

type ControllerInput struct {
	Pin              machine.Pin
	PinMode          machine.PinMode
	SwitchMode       InputType
	SendControlCodes func(bool)
}

// Configures inputs and thier corresponding MIDI outputs
// bankOffset adds a number to the each MIDI codes so multiple "Banks"
// of MIDI codes can be sent.
func initControllerInputs(bankOffset uint8) []ControllerInput {
	return []ControllerInput{
		{
			Pin:        machine.GP0,
			PinMode:    machine.PinInputPulldown,
			SwitchMode: SWITCH_HOLD,
			SendControlCodes: func(hold bool) {
				if hold {
					sendMidiCode(START_CC + 1)
				} else {
					sendMidiCode(START_CC)
				}
			},
		},

		{
			Pin:        machine.GP1,
			PinMode:    machine.PinInputPulldown,
			SwitchMode: SWITCH_HOLD,
			SendControlCodes: func(hold bool) {
				if hold {
					sendMidiCode(START_CC + 3)
				} else {
					sendMidiCode(START_CC + 2)
				}
			},
		},

		{
			Pin:        machine.GP2,
			PinMode:    machine.PinInputPulldown,
			SwitchMode: SWITCH_HOLD,
			SendControlCodes: func(hold bool) {
				if hold {
					sendMidiCode(START_CC + 5)
				} else {
					sendMidiCode(START_CC + 4)
				}
			},
		},

		{
			Pin:        machine.GP3,
			PinMode:    machine.PinInputPulldown,
			SwitchMode: SWITCH_HOLD,
			SendControlCodes: func(hold bool) {
				if hold {
					sendMidiCode(START_CC + 7)
				} else {
					sendMidiCode(START_CC + 6)
				}
			},
		},

		{
			Pin:        machine.GP4,
			PinMode:    machine.PinInputPulldown,
			SwitchMode: SWITCH_BANKSELECT,
			SendControlCodes: func(hold bool) {
				if hold {
					sendMidiCode(START_CC + 9)
				} else {
					sendMidiCode(START_CC + 8)
				}
			},
		},

		{
			Pin:        machine.GP5,
			PinMode:    machine.PinInputPulldown,
			SwitchMode: SWITCH_ONESHOT,
			SendControlCodes: func(hold bool) {
				sendMidiCode(START_CC + 10)
			},
		},
	}
}

func sendMidiCode(c uint8) {
	m := midi.New()
	OB_LED.High()
	m.SendCC(0, 0, c, MIDI_VALUE_ON)
	time.Sleep(SLEEP_5MS)
	OB_LED.Low()
}

func switchIsHeld(pin machine.Pin) bool {
	trigger := time.Now()

	for time.Now().Sub(trigger) < SWITCH_HOLD_TIME {
		if !pin.Get() {
			return false
		}
		time.Sleep(5 * time.Millisecond)
	}
	return true
}

func selectBank(currentBank uint8) uint8 {
	if currentBank == NUM_INPUT_BANKS {
		return 0
	} else {
		return currentBank + 1
	}
}

func main() {
	// Configure onboard LED
	OB_LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	// Initialize controller inputs to first bank
	inputBank := uint8(0)
	ctrlInputs := initControllerInputs(inputBank)

	for _, inp := range ctrlInputs {
		inp.Pin.Configure(machine.PinConfig{Mode: inp.PinMode})
	}

	for {
		time.Sleep(SLEEP_5MS)

		for _, inp := range ctrlInputs {

			if inp.Pin.Get() {
				switch inp.SwitchMode {
				case SWITCH_BANKSELECT:
					if switchIsHeld(inp.Pin) {
						inputBank = selectBank(inputBank)
						ctrlInputs = initControllerInputs(11 * inputBank)
					} else {
						inp.SendControlCodes(false)
					}
				case SWITCH_HOLD:
					if switchIsHeld(inp.Pin) {
						inp.SendControlCodes(true)
						time.Sleep(DEBOUNCE_TIME)
					} else {
						inp.SendControlCodes(false)
					}
				case SWITCH_ONESHOT:
					inp.SendControlCodes(false)
				}
			}
		}
	}
}
