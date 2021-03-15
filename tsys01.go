package main

import (
	//"fmt"
	"time"
	"math"
	"encoding/binary"
	"github.com/d2r2/go-i2c"
)

const (
    TSYS01_ADDR               = 0x77
    TSYS01_PROM_READ          = 0xA0
    TSYS01_RESET              = 0x1E
    TSYS01_CONVERT            = 0x48
    TSYS01_READ               = 0x00
	TSYS01_K0_ADDR = 0xAA
	TSYS01_K1_ADDR = 0xA8
	TSYS01_K2_ADDR = 0xA6
	TSYS01_K3_ADDR = 0xA4
	TSYS01_K4_ADDR = 0xA2
)

type TSYS01 struct {
	K0 float64
	K1 float64
	K2 float64
	K3 float64
	K4 float64
}

func NewTSYS01() *TSYS01 {
	d := new(TSYS01)
	return d
}

func (d *TSYS01) Init(i2c *i2c.I2C) error {

	_, err := i2c.WriteBytes([]byte{TSYS01_RESET})
	if err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)

	buff, _, err := i2c.ReadRegBytes(TSYS01_K0_ADDR, 2)
	if err != nil {
		return err
	}
	d.K0 = float64(binary.BigEndian.Uint16(buff))

	buff, _, err = i2c.ReadRegBytes(TSYS01_K1_ADDR, 2)
	if err != nil {
		return err
	}
	d.K1 = float64(binary.BigEndian.Uint16(buff))

	buff, _, err = i2c.ReadRegBytes(TSYS01_K2_ADDR, 2)
	if err != nil {
		return err
	}
	d.K2 = float64(binary.BigEndian.Uint16(buff))

	buff, _, err = i2c.ReadRegBytes(TSYS01_K3_ADDR, 2)
	if err != nil {
		return err
	}
	d.K3 = float64(binary.BigEndian.Uint16(buff))

	buff, _, err = i2c.ReadRegBytes(TSYS01_K4_ADDR, 2)
	if err != nil {
		return err
	}
	d.K4 = float64(binary.BigEndian.Uint16(buff))


	return nil
}

func (d *TSYS01) Read(i2c *i2c.I2C) (float64, error) {
	
	_, err := i2c.WriteBytes([]byte{TSYS01_CONVERT})
	if err != nil {
		return 0, err
	}
	time.Sleep(10 * time.Millisecond)

	adc_arr, _, err := i2c.ReadRegBytes(TSYS01_READ, 3)
	if err != nil {
		return 0, err
	}
	adc := float64(binary.BigEndian.Uint32([]byte{0x00, adc_arr[0], adc_arr[1], adc_arr[2]}) / 256)
	temperature := -2 * d.K4 * math.Pow10(-21) * math.Pow(adc, 4) +
				    4 * d.K3 * math.Pow10(-16) * math.Pow(adc, 3) +
				   -2 * d.K2 * math.Pow10(-11) * math.Pow(adc, 2) +
				    1 * d.K1 * math.Pow10(-6 ) * adc +
				 -1.5 * d.K0 * math.Pow10(-2 )
	return temperature, nil
}
