package main

import (
	"time"

	"machine"
	"machine/usb/adc/midi"
)

// DEBOUNCE_TIME controls how long to ignore subsequent input
const DEBOUNCE_TIME = 20 * time.Millisecond

// Define the amount of time to determine a held down switch
const SWITCH_HOLD_TIME = 500 * time.Millisecond

// Default sleep duration
const SLEEP_5MS = 5 * time.Millisecond

type SwitchAlgo interface {
	sendMidiCode(c *CtrlInput)
}

type CtrlInput struct {
	Pin         machine.Pin
	SwitchType  SwitchAlgo
	DisplayText string
}

type OneShotSwitch struct {
	midiValue   uint8
	controlCode uint8
}

type HoldSwitch struct {
	midiValueA   uint8
	controlCodeA uint8
	midiValueB   uint8
	controlCodeB uint8
}

type ProgramChangeSwitch struct {
	program uint8
}

type ToggleSwitch struct {
	state       uint8
	controlCode uint8
}

func NewCtrlInput(text string, p machine.Pin, s SwitchAlgo) *CtrlInput {
	return &CtrlInput{
		Pin:         p,
		SwitchType:  s,
		DisplayText: text,
	}
}

func (i *CtrlInput) sendControl() {
	i.SwitchType.sendMidiCode(i)
}

func (s *OneShotSwitch) sendMidiCode(c *CtrlInput) {
	m := midi.New()
	_ = m.ControlChange(0, 1, s.controlCode, s.midiValue)
}

func (s *HoldSwitch) sendMidiCode(c *CtrlInput) {
	m := midi.New()
	if switchIsHeld(c.Pin) {
		_ = m.ControlChange(0, 1, s.controlCodeB, s.midiValueB)
	} else {
		_ = m.ControlChange(0, 1, s.controlCodeA, s.midiValueA)
	}
}

func (s *ProgramChangeSwitch) sendMidiCode(c *CtrlInput) {
	m := midi.New()
	_ = m.ControlChange(0, 1, 32, 0)
	_ = m.ProgramChange(0, 1, s.program)
}

func (s *ToggleSwitch) sendMidiCode(c *CtrlInput) {
	m := midi.New()
	if s.state == MIDI_VALUE_ON {
		s.state = MIDI_VALUE_OFF
	} else {
		s.state = MIDI_VALUE_ON
	}
	_ = m.ControlChange(0, 1, s.controlCode, s.state)

}

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

func waitForSwitchRelease(pin machine.Pin) {
	for {
		if !pin.Get() {
			time.Sleep(DEBOUNCE_TIME)
			return
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
			waitForSwitchRelease(SWITCHPINS[5])
		}

		for _, inp := range ctrlInputs[inputBank] {

			if inp.Pin.Get() {
				inp.sendControl()
				waitForSwitchRelease(inp.Pin)
			}
		}
	}
}
