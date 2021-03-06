package controller

import (
	"fmt"
	"time"
)

type Strip struct {
	Name  string `json:"name"`
	Color HSI    `json:"hsi"`
	rPin  pwmPin
	gPin  pwmPin
	bPin  pwmPin
}

func NewStrip(name string, rPinNumber, gPinNumber, bPinNumber string) (*Strip, error) {

	rPin, err := newPWMPin(rPinNumber)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize rPin: %v", err)
	}

	gPin, err := newPWMPin(gPinNumber)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize gPin: %v", err)
	}

	bPin, err := newPWMPin(bPinNumber)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize bPin: %v", err)
	}

	s := Strip{
		Name: name,
		rPin: *rPin,
		gPin: *gPin,
		bPin: *bPin,
	}

	return &s, nil

}

func (s *Strip) SetColor(color HSI) error {
	s.Color = color
	return s.setPins()
}

func (s *Strip) setPins() error {

	color := s.Color.ToRGB()

	if err := s.rPin.Set(color.Red); err != nil {
		return fmt.Errorf("Failed to set rPin: %v", err)
	}

	if err := s.gPin.Set(color.Green); err != nil {
		return fmt.Errorf("Failed to set gPin: %v", err)
	}

	if err := s.bPin.Set(color.Blue); err != nil {
		return fmt.Errorf("Failed to set bPin: %v", err)
	}
	return nil
}

func (s *Strip) TestStrip() error {

	const testSeparationDuration = 250

	fmt.Printf("Testing LED Strip %s", s.Name)

	var test TestPatterns
	test.Default()

	for _, v := range test {
		fmt.Printf("Starting Test %s\n", v.Name)
		err := s.SetColor(v.Color)
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(v.Duration) * time.Millisecond)
		err = s.Off()
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(testSeparationDuration) * time.Millisecond)

	}
	return nil
}

func (s *Strip) Stop() {
	ps.Pub(true, s.Name)
}

func (s *Strip) StopChan() chan interface{} {
	return ps.Sub(s.Name)
}

func (s *Strip) Unsub(ch chan interface{}) {
	ps.Unsub(ch, s.Name)
}
