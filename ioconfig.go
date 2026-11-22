package main

import "machine"

const BANKCFGINP uint8 = 5
const NUMINPUTS = 6

// Common value assingment MIDI ON
const MIDI_VALUE_ON uint8 = 0x7f
const MIDI_VALUE_OFF uint8 = 0x00

var SWITCHPINS []machine.Pin = []machine.Pin{
	machine.GP0, machine.GP1, machine.GP2, machine.GP3, machine.GP4, machine.GP5,
}

func ConfigureSwitchPins() {
	for _, pin := range SWITCHPINS {
		pin.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}
}

type SwitchBank []*CtrlInput

// If no EEPROM configuration Found or config is bad,
// These defaults will be used.
func InitDefaultSwitchBanks() []SwitchBank {
	return []SwitchBank{
		// Bank 0: Program/Preset Change
		{
			NewCtrlInput("P1", machine.GP0, &ProgramChangeSwitch{0}),
			NewCtrlInput("P2", machine.GP1, &ProgramChangeSwitch{1}),
			NewCtrlInput("P3", machine.GP2, &ProgramChangeSwitch{2}),
			NewCtrlInput("P4", machine.GP3, &ProgramChangeSwitch{3}),
			NewCtrlInput("P5", machine.GP4, &ProgramChangeSwitch{4}),
		},
		// Bank 1: Scene Change
		{
			NewCtrlInput("1-2", machine.GP0, &HoldSwitch{
				midiValueA:   0x00,
				controlCodeA: 69,
				midiValueB:   0x01,
				controlCodeB: 69,
			}),
			NewCtrlInput("3-4", machine.GP1, &HoldSwitch{
				midiValueA:   0x02,
				controlCodeA: 69,
				midiValueB:   0x03,
				controlCodeB: 69,
			}),
			NewCtrlInput("5-6", machine.GP2, &HoldSwitch{
				midiValueA:   0x04,
				controlCodeA: 69,
				midiValueB:   0x05,
				controlCodeB: 69,
			}),
			NewCtrlInput("7-8", machine.GP3, &HoldSwitch{
				midiValueA:   0x06,
				controlCodeA: 69,
				midiValueB:   0x07,
				controlCodeB: 69,
			}),
			NewCtrlInput("TAP", machine.GP4, &OneShotSwitch{
				midiValue:   MIDI_VALUE_ON,
				controlCode: 64,
			}),
		},
		// Bank 2: Effects Toggle
		{
			NewCtrlInput("OD", machine.GP0, &ToggleSwitch{
				state:       MIDI_VALUE_OFF,
				controlCode: 33,
			}),
			NewCtrlInput("DLAY", machine.GP1, &ToggleSwitch{
				state:       MIDI_VALUE_OFF,
				controlCode: 34,
			}),
			NewCtrlInput("DIVE", machine.GP2, &ToggleSwitch{
				state:       MIDI_VALUE_OFF,
				controlCode: 35,
			}),
			NewCtrlInput("MOD", machine.GP3, &ToggleSwitch{
				state:       MIDI_VALUE_OFF,
				controlCode: 36,
			}),
			NewCtrlInput("BST", machine.GP4, &ToggleSwitch{
				state:       MIDI_VALUE_OFF,
				controlCode: 37,
			}),
		},
		// Bank 3: Looper
		{
			NewCtrlInput("REC", machine.GP0, &OneShotSwitch{
				midiValue:   MIDI_VALUE_ON,
				controlCode: 60,
			}),
			NewCtrlInput("ODUB", machine.GP1, &OneShotSwitch{
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 60,
			}),
			NewCtrlInput("PLAY", machine.GP2, &HoldSwitch{
				midiValueA:   MIDI_VALUE_ON,
				controlCodeA: 61,
				midiValueB:   MIDI_VALUE_OFF,
				controlCodeB: 61,
			}),
			NewCtrlInput("UNDO", machine.GP3, &OneShotSwitch{
				midiValue:   MIDI_VALUE_ON,
				controlCode: 63,
			}),
			NewCtrlInput("TAP", machine.GP4, &OneShotSwitch{
				midiValue:   MIDI_VALUE_ON,
				controlCode: 64,
			}),
		},
	}
}
