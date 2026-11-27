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
const SLEEP_TIME = 3 * time.Millisecond

// Watchdog timeout – long enough for normal hold detection + margin,
// short enough that a real hang recovers quickly.
const WATCHDOG_TIMEOUT_MS = 2500

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

var midiPort = midi.Port()

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
	_ = midiPort.ControlChange(0, 1, s.controlCode, s.midiValue)
}

func (s *HoldSwitch) sendMidiCode(c *CtrlInput) {
	if switchIsHeld(c.Pin) {
		_ = midiPort.ControlChange(0, 1, s.controlCodeB, s.midiValueB)
	} else {
		_ = midiPort.ControlChange(0, 1, s.controlCodeA, s.midiValueA)
	}
}

func (s *ProgramChangeSwitch) sendMidiCode(c *CtrlInput) {
	_ = midiPort.ControlChange(0, 1, 32, 0)
	_ = midiPort.ProgramChange(0, 1, s.program)
}

func (s *ToggleSwitch) sendMidiCode(c *CtrlInput) {
	if s.state == MIDI_VALUE_ON {
		s.state = MIDI_VALUE_OFF
	} else {
		s.state = MIDI_VALUE_ON
	}
	_ = midiPort.ControlChange(0, 1, s.controlCode, s.state)
}

func switchIsHeld(pin machine.Pin) bool {
	trigger := time.Now()

	for time.Now().Sub(trigger) < SWITCH_HOLD_TIME {
		machine.Watchdog.Update() // keep the dog happy while we wait
		if !pin.Get() {
			return false
		}
		time.Sleep(SLEEP_TIME)
	}
	return true
}

func selectBank(numOfBanks uint8, currentBank uint8, reverse bool) uint8 {
	if reverse {
		if currentBank == 0 {
			return numOfBanks - 1
		}
		return currentBank - 1
	}
	if currentBank == numOfBanks-1 {
		return 0
	}
	return currentBank + 1
}

func waitForSwitchRelease(pin machine.Pin) {
	for pin.Get() {
		machine.Watchdog.Update()
		time.Sleep(SLEEP_TIME)
	}
	time.Sleep(DEBOUNCE_TIME)
}

func main() {
	// Configure & start the hardware watchdog
	err := machine.Watchdog.Configure(machine.WatchdogConfig{
		TimeoutMillis: WATCHDOG_TIMEOUT_MS,
	})
	if err != nil {
		// Very unlikely on RP2040, but good practice
		println("watchdog configure failed:", err)
	}
	err = machine.Watchdog.Start()
	if err != nil {
		println("watchdog start failed:", err)
	}

	disp := Display{initDisplay()}

	inputBank := uint8(0)

	ConfigureSwitchPins()

	ctrlInputs := InitDefaultSwitchBanks()
	numOfBanks := uint8(len(ctrlInputs))
	bankSelect := SWITCHPINS[5]

	disp.WriteOut([]byte(formatDisplayText(ctrlInputs[inputBank], inputBank)))

	for {
		machine.Watchdog.Update()

		time.Sleep(SLEEP_TIME)

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
