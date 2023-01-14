package main

import (
	"fmt"
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

// Control code assignments
// The code range 20-31 is used because it is undefined by the MIDI standard
const (
	// Next Track
	MIDI_CC20 uint8 = iota + 0x14
	// Previous Track
	MIDI_CC21
	// Mute Track
	MIDI_CC22
	// Solo Track
	MIDI_CC23
	// Tap Tempo
	MIDI_CC24
	// Unassigned
	MIDI_CC25
	// Play
	MIDI_CC26
	// Stop
	MIDI_CC27
	// Record
	MIDI_CC28
	// Delete Recording
	MIDI_CC29
	// Arm Recording
	MIDI_CC30
	// Unassigned
	MIDI_CC31
)

type ControlCodes uint8

// Common value assingments for control codes
const (
	// MIDI_VALUE_OFF uint8 = 0x00
	MIDI_VALUE_ON uint8 = 0x7f
)

const (
	SWITCH_ONESHOT uint = iota
	SIWTCH_TWOSHOT
	SWITCH_HOLD
)

type InputType uint8

type ControllerInput struct {
	Pin              machine.Pin
	PinMode          machine.PinMode
	SwitchMode       InputType
	SendControlCodes func(bool)
}

const PINS_IN_USE = 6

func initControllerInputs() [PINS_IN_USE]ControllerInput {

	var tksel ControllerInput = ControllerInput{
		Pin:        machine.GP0,
		PinMode:    machine.PinInputPulldown,
		SwitchMode: InputType(SWITCH_HOLD),
		SendControlCodes: func(hold bool) {
			if hold {
				sendMidiCode([]uint8{MIDI_CC21})
			} else {
				sendMidiCode([]uint8{MIDI_CC20})
			}
		},
	}

	var muteSolo ControllerInput = ControllerInput{
		Pin:        machine.GP1,
		PinMode:    machine.PinInputPulldown,
		SwitchMode: InputType(SWITCH_HOLD),
		SendControlCodes: func(hold bool) {
			if hold {
				sendMidiCode([]uint8{MIDI_CC23})
			} else {
				sendMidiCode([]uint8{MIDI_CC22})
			}
		},
	}

	var tapTempo ControllerInput = ControllerInput{
		Pin:        machine.GP2,
		PinMode:    machine.PinInputPulldown,
		SwitchMode: InputType(SWITCH_ONESHOT),
		SendControlCodes: func(hold bool) {
			fmt.Println("Sending CC24")
			sendMidiCode([]uint8{MIDI_CC24})
			fmt.Println("Sent CC24")
		},
	}

	var playstop ControllerInput = ControllerInput{
		Pin:        machine.GP3,
		PinMode:    machine.PinInputPulldown,
		SwitchMode: InputType(SWITCH_HOLD),
		SendControlCodes: func(hold bool) {
			if hold {
				sendMidiCode([]uint8{MIDI_CC27})
			} else {
				sendMidiCode([]uint8{MIDI_CC26})
			}
		},
	}

	var record ControllerInput = ControllerInput{
		Pin:        machine.GP4,
		PinMode:    machine.PinInputPulldown,
		SwitchMode: InputType(SWITCH_HOLD),
		SendControlCodes: func(hold bool) {
			if hold {
				sendMidiCode([]uint8{MIDI_CC29})
			} else {
				sendMidiCode([]uint8{MIDI_CC28})
			}
		},
	}

	var armRecord ControllerInput = ControllerInput{
		Pin:        machine.GP5,
		PinMode:    machine.PinInputPulldown,
		SwitchMode: InputType(SWITCH_ONESHOT),
		SendControlCodes: func(hold bool) {
			sendMidiCode([]uint8{MIDI_CC30})
		},
	}

	return [PINS_IN_USE]ControllerInput{
		tksel, muteSolo, tapTempo, playstop, record, armRecord,
	}
}

func sendMidiCode(c []uint8) {
	m := midi.New()
	for _, code := range c {
		m.SendCC(0, 0, code, MIDI_VALUE_ON)
		time.Sleep(SLEEP_5MS)
	}

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

func main() {
	// Initialize onboard LED
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ctrlInputs := initControllerInputs()

	for _, inp := range ctrlInputs {
		inp.Pin.Configure(machine.PinConfig{Mode: inp.PinMode})
	}

	for {
		time.Sleep(SLEEP_5MS)

		for _, inp := range ctrlInputs {

			if inp.Pin.Get() {
				switch inp.SwitchMode {

				case InputType(SWITCH_HOLD):
					if switchIsHeld(inp.Pin) {
						led.High()
						inp.SendControlCodes(true)
						time.Sleep(DEBOUNCE_TIME)
						led.Low()
					} else {
						led.High()
						inp.SendControlCodes(false)
						led.Low()
					}

				case InputType(SWITCH_ONESHOT):
					led.High()
					inp.SendControlCodes(false)
					led.Low()
				}
			}
		}
	}
}
