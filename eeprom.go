package main

import (
	"encoding/binary"
	"fmt"
	"machine"

	"tinygo.org/x/drivers/at24cx"
)

const (
	// Byte index for config length
	EEPROM_IDX_CFG_LEN = 0
	// Byte index for config CRC
	EEPROM_IDX_CFG_CRC = 1
	// Byte index for config
	EEPROM_IDX_CFG = 5
)

type EEPROM struct {
	*at24cx.Device
}

func initEEPROM() *at24cx.Device {
	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{
		SCL: machine.GP13,
		SDA: machine.GP12,
	})
	if err != nil {
		fmt.Println(err)
	}
	d := at24cx.New(i2c)
	d.Address = 0x50
	d.Configure(at24cx.Config{
		PageSize:        64,
		StartRAMAddress: 0,
		EndRAMAddress:   4096,
	})
	return &d
}

func (d *EEPROM) getStoredConfigLen() (uint16, error) {
	// sz, err := d.ReadAt([]byte{0, 1}, 2)
	// if err != nil {
	// 	fmt.Println("Read error: ", err)
	// }
	// return sz
	data := make([]byte, 2)
	_, err := d.ReadAt(data, EEPROM_IDX_CFG_LEN)
	if err != nil {
		fmt.Println("Read error: ", err)
	}
	return binary.BigEndian.Uint16(data), nil
}

func (d *EEPROM) setStoredConfigLen(size uint16) error {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, size)
	_, err := d.WriteAt(data, EEPROM_IDX_CFG_LEN)
	if err != nil {
		return err
	}
	return nil
}

func (d *EEPROM) getStoredConfig(size uint16) ([]byte, error) {
	data := make([]byte, size)
	// startByte := uint16(EEPROM_IDX_CFG)
	// for i := startByte; i < size; i++ {
	// 	b, err := d.ReadByte(i)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	data = append(data, b)
	// }
	_, err := d.ReadAt(data, EEPROM_IDX_CFG)
	if err != nil {
		return nil, err
	}
	return data, nil
}
