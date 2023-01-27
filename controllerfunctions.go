package main

import (
	"machine"
	"machine/usb/midi"
)

type SwitchAlgo interface {
	sendMidiCode(c *CtrlInput)
}

type CtrlInput struct {
	Pin         machine.Pin
	SwitchType  SwitchAlgo
	DisplayText string
}

func NewCtrlInput(text string, p machine.Pin, s SwitchAlgo) *CtrlInput {
	return &CtrlInput{
		Pin:         p,
		SwitchType:  s,
		DisplayText: text,
	}
}

func selectMidiValue(toggle bool, state uint8) uint8 {
	if !toggle {
		return MIDI_VALUE_ON
	} else {
		return state ^ MIDI_VALUE_ON
	}
}

func (i *CtrlInput) sendControl() {
	i.SwitchType.sendMidiCode(i)
}

type OneShotSwitch struct {
	toggle      bool
	midiValue   uint8
	controlCode uint8
}

func (s *OneShotSwitch) sendMidiCode(c *CtrlInput) {
	s.midiValue = selectMidiValue(s.toggle, s.midiValue)
	m := midi.New()
	m.SendCC(0, 0, s.controlCode, s.midiValue)
}

type HoldSwitch struct {
	toggleA      bool
	midiValueA   uint8
	controlCodeA uint8
	toggleB      bool
	midiValueB   uint8
	controlCodeB uint8
}

func (s *HoldSwitch) sendMidiCode(c *CtrlInput) {
	m := midi.New()
	if switchIsHeld(c.Pin) {
		s.midiValueB = selectMidiValue(s.toggleB, s.midiValueB)
		m.SendCC(0, 0, s.controlCodeB, s.midiValueB)
	} else {
		s.midiValueA = selectMidiValue(s.toggleA, s.midiValueA)
		m.SendCC(0, 0, s.controlCodeA, s.midiValueA)
	}
}
