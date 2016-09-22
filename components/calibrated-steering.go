package components

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"

	"sdmimaye.de/smart-video-car/stream"

	"github.com/go-ini/ini"
)

const (
	servoIndex = 0
)

//CalibratedSteering controls the vehicle steering. It composes of one servo motors and the required configuration
type CalibratedSteering struct {
	servo CalibratedServo
	left  float64
	right float64
}

func doLoadIniWithMatchingSectionOrCreateEmptyForSteering() (*ini.File, *ini.Section, error) {
	cfg, err := ini.LooseLoad(CalibrationFilePath)
	if err != nil {
		return nil, nil, errors.New("Could not Load ini file. Error: " + err.Error())
	}

	sec, err := cfg.GetSection("Control")
	if err != nil {
		log.Print("Missing Section: Control. Will regenerate section\n")
		sec, err = cfg.NewSection("Control")
		if err != nil {
			return nil, nil, errors.New("Error while generating Section. Error: " + err.Error())
		}
	}
	cfg.SaveTo(CalibrationFilePath)

	return cfg, sec, nil
}

//NewCalibratedSteering will create a new calibrated vehicle steering
func NewCalibratedSteering() (*CalibratedSteering, error) {
	ctrl := CalibratedSteering{}
	servo, err := NewCalibratedServo(servoIndex)
	if err != nil {
		return nil, fmt.Errorf("Steering: Could not create calibrated servo on channel: %v. Error: %v", servoIndex, err)
	}
	ctrl.servo = *servo
	_, section, err := doLoadIniWithMatchingSectionOrCreateEmptyForSteering()
	if err != nil {
		return nil, fmt.Errorf("Could not create calibrated steering. Error: %v", err)
	}

	ctrl.left, _ = section.Key("Left").Float64()
	ctrl.right, _ = section.Key("Right").Float64()

	log.Printf("Current Steering-Config: %v\n", ctrl)
	return &ctrl, nil
}

//Calibrate will calibrate the steering
func (c *CalibratedSteering) Calibrate(stream stream.Stream) error {
	servo, err := NewCalibratedServo(servoIndex)
	if err != nil {
		return fmt.Errorf("Steering: Could not create calibrated servo on channel: %v. Error: %v", servoIndex, err)
	}

	err = servo.Calibrate(stream)
	if err != nil {
		return fmt.Errorf("Steering: Could not calibrated servo on channel: %v. Error: %v", servoIndex, err)
	}

	r := stream.GetReader()
	w := stream.GetWriter()

	reader := bufio.NewReader(r)
	for {
		fmt.Fprint(w, "Steering will move to home position. Press any key to continue...\r\n")
		err = servo.Home()
		if err != nil {
			return fmt.Errorf("Could not determine Steering-Calibration. Error while moving to home position: %v", err)
		}
		reader.ReadString('\n')

		fmt.Fprint(w, "Steering will move in positive direction. Please pay attention to the direction! Press any key to continue...\r\n")
		reader.ReadString('\n')

		err = servo.Forward(100)
		if err != nil {
			return fmt.Errorf("Could not determine Steering-Calibration. Error while moving forward: %v", err)
		}
		fmt.Fprint(w, "Did the steering move [L]eft or [R]ight (press anything else to repeat)?\r\n")
		val, _ := reader.ReadString('\n')
		if strings.HasPrefix(val, "L") {
			c.left = 1.0
			c.right = -1.0
			break
		} else if strings.HasPrefix(val, "R") {
			c.left = -1.0
			c.right = 1.0
			break
		}
	}

	cfg, section, err := doLoadIniWithMatchingSectionOrCreateEmptyForSteering()
	if err != nil {
		return fmt.Errorf("Could not create config file for steering. Error: %v", err)
	}

	section.NewKey("Left", fmt.Sprintf("%v", c.left))
	section.NewKey("Right", fmt.Sprintf("%v", c.right))

	err = cfg.SaveTo(CalibrationFilePath)
	if err != nil {
		return fmt.Errorf("Could not store configuration for steering. Error: %v", err)
	}

	return nil
}

//SteerLeft will steer the vehicle in a left position
func (c *CalibratedSteering) SteerLeft(percent float64) error {
	return c.servo.Move(percent, c.left)
}

//SteerRight will steer the vehicle in a right position
func (c *CalibratedSteering) SteerRight(percent float64) error {
	return c.servo.Move(percent, c.right)
}
