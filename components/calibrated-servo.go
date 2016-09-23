package components

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"sdmimaye.de/smart-video-car/hardware"
	"sdmimaye.de/smart-video-car/stream"

	"github.com/go-ini/ini"
)

const (
	//ServoCalibrationPrefix is the prefix for the different servo motors
	ServoCalibrationPrefix = "Servo"
)

//CalibratedServo represents a calibrateable servo in a smart car
type CalibratedServo struct {
	channel int
	servo   hardware.ServoMotor
	min     int
	max     int
	center  int
}

//Home will move the servo in a centered direction (percentual to the maximum calibrated value)
func (s *CalibratedServo) Home() error {
	return s.servo.SetAngle(s.center)
}

func doCalculatePercentOfAndSteer(s *CalibratedServo, percent float64) error {
	min := float64(s.min)
	max := float64(s.max)
	center := float64(s.center)
	value := 0.0

	if percent >= 0 && percent <= 100 { //positive movement
		value = center + ((max - center) / 100 * percent)
	} else if percent < 0 && percent >= -100 { //negative movement
		value = min + ((center - min) / 100 * percent)
	} else {
		return errors.New("Invalid percentual value for Servo")
	}

	return s.servo.SetAngle(int(value))
}

//Forward will move the servo in a forward direction (percentual to the maximum calibrated value)
func (s *CalibratedServo) Forward(percent float64) error {
	return doCalculatePercentOfAndSteer(s, percent)
}

//Backward will move the servo in a backward direction (percentual to the maximum calibrated value)
func (s *CalibratedServo) Backward(percent float64) error {
	return doCalculatePercentOfAndSteer(s, -percent)
}

//Move will move the servo in a direction (based on the passed direction)/(percentual to the maximum calibrated value)
func (s *CalibratedServo) Move(percent float64, direction float64) error {
	if direction < 0 {
		direction = -1
	} else {
		direction = 1
	}

	return doCalculatePercentOfAndSteer(s, percent*direction)
}

func doDetermineAngle(r io.Reader, w io.Writer, servo hardware.ServoMotor, angle int, direction string) int {
	reader := bufio.NewReader(r)

	for {
		fmt.Fprintf(w, "Current %v Value: %v\r\n", direction, angle)
		servo.SetAngle(angle)
		fmt.Fprintf(w, "%v Calibration. Press [I] for a increment (10), [i] for a small increment (1), [D] for a big decrement (10), [d] for a small decrement (1) and [X] to end calibration...\r\n", direction)
		val, _ := reader.ReadString('\n')

		if strings.HasPrefix(val, "I") {
			angle += 10
		} else if strings.HasPrefix(val, "i") {
			angle++
		} else if strings.HasPrefix(val, "D") {
			angle -= 10
		} else if strings.HasPrefix(val, "d") {
			angle--
		} else {
			fmt.Fprintf(w, "Setting %v for %v value\r\n", angle, direction)
			return angle
		}
	}
}

func doLoadIniWithMatchingSectionOrCreateEmptyForServo(section int) (*ini.File, *ini.Section, error) {
	cfg, err := ini.LooseLoad(CalibrationFilePath)
	if err != nil {
		return nil, nil, errors.New("Could not Load ini file. Error: " + err.Error())
	}

	name := fmt.Sprintf("%v%v", ServoCalibrationPrefix, section)
	sec, err := cfg.GetSection(name)
	if err != nil {
		log.Printf("Missing Section: %v. Will regenerate section\n", name)
		sec, err = cfg.NewSection(name)
		if err != nil {
			return nil, nil, errors.New("Error while generating Section. Error: " + err.Error())
		}
	}

	if !sec.Haskey("Min") {
		_, err = sec.NewKey("Min", "90")
		if err != nil {
			return nil, nil, errors.New("Could not generate Min-Value for Section. Error: " + err.Error())
		}
	}

	if !sec.Haskey("Max") {
		_, err = sec.NewKey("Max", "90")
		if err != nil {
			return nil, nil, errors.New("Could not generate Max-Value for Section. Error: " + err.Error())
		}
	}

	if !sec.Haskey("Center") {
		_, err = sec.NewKey("Center", "90")
		if err != nil {
			return nil, nil, errors.New("Could not generate Center-Value for Section. Error: " + err.Error())
		}
	}
	cfg.SaveTo(CalibrationFilePath)

	return cfg, sec, nil
}

//Calibrate will (re)calibrate a servo
func (s *CalibratedServo) Calibrate(stream stream.Stream) error {
	r := stream.GetReader()
	w := stream.GetWriter()

	fmt.Fprintf(w, "Starting calibration for servo: %v\r\n", s.channel)
	s.min = doDetermineAngle(r, w, s.servo, s.min, "Min")
	s.max = doDetermineAngle(r, w, s.servo, s.max, "Max")
	s.center = doDetermineAngle(r, w, s.servo, s.center, "Center")

	cfg, section, err := doLoadIniWithMatchingSectionOrCreateEmptyForServo(s.channel)
	if err != nil {
		return fmt.Errorf("Could not read config file for servo: %v. Error: %v", s.channel, err)
	}

	section.NewKey("Min", strconv.Itoa(s.min))
	section.NewKey("Max", strconv.Itoa(s.max))
	section.NewKey("Center", strconv.Itoa(s.center))

	err = cfg.SaveTo(CalibrationFilePath)
	if err != nil {
		return fmt.Errorf("Could not store configuration for servo: %v. Error: %v", s.channel, err)
	}
	s.servo.SetAngle(s.center)

	return nil
}

//NewCalibratedServo will create a new calibrated servo. If no calibration file is present a calibration will be initiated
func NewCalibratedServo(channel int) (*CalibratedServo, error) {
	servo, err := hardware.GetServo(channel)

	if err != nil {
		return nil, fmt.Errorf("Error while generating servo motor on channel: %v. Error: %v", channel, err)
	}

	_, sec, err := doLoadIniWithMatchingSectionOrCreateEmptyForServo(channel)
	if err != nil {
		return nil, fmt.Errorf("Error while loading ini file for servo on channel: %v, Error: %v", channel, err)
	}

	srv := CalibratedServo{channel: channel, servo: servo}
	srv.min, _ = sec.Key("Min").Int()
	srv.max, _ = sec.Key("Max").Int()
	srv.center, _ = sec.Key("Center").Int()

	return &srv, nil
}
