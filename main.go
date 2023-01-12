package main

import (
	"fmt"
	"time"

	"machine"
	"machine/usb/midi"
)

// DEBOUNCE_TIME controls how long to ignore subsequent input
const DEBOUNCE_TIME = time.Millisecond * 500

// Define the amount of time to determine a held down switch
const SWITCH_HOLD_TIME = time.Millisecond * 500

// Control code assignments
// The code range 20-31 is used because it is undefined by the MIDI standard
const (
	// Next Track
	MIDI_CC20 uint8 = iota + 0x14
	// Previous Track
	MIDI_CC21
	// Arm Recording
	MIDI_CC22
	// Play
	MIDI_CC23
	// Stop
	MIDI_CC24
	MIDI_CC25
	MIDI_CC26
	MIDI_CC27
	MIDI_CC28
	MIDI_CC29
	MIDI_CC30
	MIDI_CC31
)

type MidiControlCode uint8

// Common value assingments for control codes
const (
	MIDI_VALUE_OFF uint8 = 0x00
	MIDI_VALUE_ON  uint8 = 0x7f
	// MIDPOINT_OFF uint8 = 0x3f
	// MIDPOINT_ON  uint8 = 0x40
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

var tksel ControllerInput = ControllerInput{
	Pin:        machine.GP0,
	PinMode:    machine.PinInputPulldown,
	SwitchMode: InputType(SWITCH_HOLD),
	SendControlCodes: func(hold bool) {
		m := midi.New()
		if hold {
			m.SendCC(0, 0, MIDI_CC22, MIDI_VALUE_ON)
			time.Sleep(5 * time.Millisecond)
			m.SendCC(0, 0, MIDI_CC21, MIDI_VALUE_ON)
		} else {
			m.SendCC(0, 0, MIDI_CC22, MIDI_VALUE_ON)
			time.Sleep(5 * time.Millisecond)
			m.SendCC(0, 0, MIDI_CC20, MIDI_VALUE_ON)
		}
	},
}

var playstop ControllerInput = ControllerInput{
	Pin:        machine.GP1,
	PinMode:    machine.PinInputPulldown,
	SwitchMode: InputType(SWITCH_HOLD),
	SendControlCodes: func(hold bool) {
		m := midi.New()
		if hold {
			m.SendCC(0, 0, MIDI_CC24, MIDI_VALUE_ON)
		} else {
			m.SendCC(0, 0, MIDI_CC23, MIDI_VALUE_ON)
		}
	},
}

const PINS_IN_USE = 2

var CtrlInputs [PINS_IN_USE]ControllerInput = [PINS_IN_USE]ControllerInput{tksel, playstop}

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

	for _, inp := range CtrlInputs {
		inp.Pin.Configure(machine.PinConfig{Mode: inp.PinMode})
	}

	for {
		time.Sleep(5 * time.Millisecond)

		for _, inp := range CtrlInputs {

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
					fmt.Println("oneshot")
				}
			}
		}
	}

}
