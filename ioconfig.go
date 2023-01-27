package main

import "machine"

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
		// Bank 0
		{
			NewCtrlInput("PLAY", machine.GP0, &HoldSwitch{
				toggleA:      false,
				midiValueA:   MIDI_VALUE_OFF,
				controlCodeA: 20,
				toggleB:      false,
				midiValueB:   MIDI_VALUE_OFF,
				controlCodeB: 21,
			}),
			NewCtrlInput("TKSEL", machine.GP1, &HoldSwitch{
				toggleA:      false,
				midiValueA:   MIDI_VALUE_OFF,
				controlCodeA: 22,
				toggleB:      false,
				midiValueB:   MIDI_VALUE_OFF,
				controlCodeB: 23,
			}),
			NewCtrlInput("ARM", machine.GP2, &HoldSwitch{
				toggleA:      false,
				midiValueA:   MIDI_VALUE_OFF,
				controlCodeA: 24,
				toggleB:      false,
				midiValueB:   MIDI_VALUE_OFF,
				controlCodeB: 25,
			}),
			NewCtrlInput("REC", machine.GP3, &HoldSwitch{
				toggleA:      false,
				midiValueA:   MIDI_VALUE_OFF,
				controlCodeA: 24,
				toggleB:      false,
				midiValueB:   MIDI_VALUE_OFF,
				controlCodeB: 25,
			}),
			NewCtrlInput("TAP", machine.GP4, &OneShotSwitch{
				toggle:      false,
				midiValue:   MIDI_VALUE_ON,
				controlCode: 28,
			}),
		},
		// Bank 1
		{
			NewCtrlInput("DRIVE", machine.GP0, &OneShotSwitch{
				toggle:      true,
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 30,
			}),
			NewCtrlInput("DELAY", machine.GP1, &OneShotSwitch{
				toggle:      true,
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 32,
			}),
			NewCtrlInput("BOOST", machine.GP2, &OneShotSwitch{
				toggle:      true,
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 34,
			}),
			NewCtrlInput("CHRUS", machine.GP3, &OneShotSwitch{
				toggle:      true,
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 36,
			}),
			NewCtrlInput("COMP", machine.GP0, &OneShotSwitch{
				toggle:      true,
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 38,
			}),
		},
		// Bank 2
		{
			NewCtrlInput("N/A", machine.GP0, &OneShotSwitch{
				toggle:      true,
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 40,
			}),
			NewCtrlInput("N/A", machine.GP1, &OneShotSwitch{
				toggle:      true,
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 42,
			}),
			NewCtrlInput("N/A", machine.GP2, &OneShotSwitch{
				toggle:      true,
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 44,
			}),
			NewCtrlInput("N/A", machine.GP3, &OneShotSwitch{
				toggle:      true,
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 46,
			}),
			NewCtrlInput("N/A", machine.GP0, &OneShotSwitch{
				toggle:      true,
				midiValue:   MIDI_VALUE_OFF,
				controlCode: 48,
			}),
		},
	}
}
