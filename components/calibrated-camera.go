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

var servos = []int{14, 15}

//CameraServoConfig determines how the camera servo is configured
type CameraServoConfig struct {
	index int
	sign  float64
}

//CalibratedCamera controls the camera. It composes of two servo motors and the handle to communicate with the camera
type CalibratedCamera struct {
	servos []CalibratedServo
	up     CameraServoConfig
	down   CameraServoConfig
	left   CameraServoConfig
	right  CameraServoConfig
}

func doLoadIniWithMatchingSectionOrCreateEmptyForCamera() (*ini.File, *ini.Section, error) {
	cfg, err := ini.LooseLoad(CalibrationFilePath)
	if err != nil {
		return nil, nil, errors.New("Could not Load ini file. Error: " + err.Error())
	}

	sec, err := cfg.GetSection("Camera")
	if err != nil {
		log.Print("Missing Section: Camera. Will regenerate section\n")
		sec, err = cfg.NewSection("Camera")
		if err != nil {
			return nil, nil, errors.New("Error while generating Section. Error: " + err.Error())
		}
	}
	cfg.SaveTo(CalibrationFilePath)

	return cfg, sec, nil
}

//NewCalibratedCamera will create a new calibrated camera
func NewCalibratedCamera() (*CalibratedCamera, error) {
	cam := CalibratedCamera{}

	cam.servos = make([]CalibratedServo, len(servos))
	for i, channel := range servos {
		servo, err := NewCalibratedServo(channel)
		if err != nil {
			return nil, fmt.Errorf("Could not create calibrated servo on channel: %v. Error: %v", channel, err)
		}

		cam.servos[i] = *servo
	}

	_, section, err := doLoadIniWithMatchingSectionOrCreateEmptyForCamera()
	if err != nil {
		return nil, fmt.Errorf("Could not create calibrated Camera. Error: %v", err)
	}

	cam.up.index, _ = section.Key("CameraUpIndex").Int()
	cam.up.sign, _ = section.Key("CameraUpSign").Float64()

	cam.down.index, _ = section.Key("CameraDownIndex").Int()
	cam.down.sign, _ = section.Key("CameraDownSign").Float64()

	cam.left.index, _ = section.Key("CameraLeftIndex").Int()
	cam.left.sign, _ = section.Key("CameraLeftSign").Float64()

	cam.right.index, _ = section.Key("CameraRightIndex").Int()
	cam.right.sign, _ = section.Key("CameraRightSign").Float64()

	log.Printf("Current Camera-Config: %v\n", cam)
	return &cam, nil
}

//Calibrate will calibrate the camera
func (s *CalibratedCamera) Calibrate(stream stream.Stream) error {
	s.servos = make([]CalibratedServo, len(servos))
	for i, channel := range servos {
		servo, err := NewCalibratedServo(channel)
		if err != nil {
			return fmt.Errorf("Could not create calibrated servo on channel: %v. Error: %v", channel, err)
		}

		err = servo.Calibrate(stream)
		if err != nil {
			return fmt.Errorf("Camera: Could not calibrate servo on channel: %v for camera. Error: %v", channel, err)
		}

		r := stream.GetReader()
		w := stream.GetWriter()

		reader := bufio.NewReader(r)
		for {
			fmt.Fprint(w, "Camera will move to home position. Press any key to continue...\r\n")
			err = servo.Home()
			if err != nil {
				return fmt.Errorf("Could not determine Camera-Calibration. Error while moving to home position: %v\r\n", err)
			}
			reader.ReadString('\n')

			fmt.Fprint(w, "Camera will move in positive direction. Please pay attention to the direction! Press any key to continue...\r\n")
			reader.ReadString('\n')

			err = servo.Forward(100)
			if err != nil {
				return fmt.Errorf("Could not determine Camera-Calibration. Error while moving forward: %v", err)
			}
			fmt.Fprint(w, "Did the camera move [U]p, [D]own, [L]eft or [R]ight (press anything else to repeat)?\r\n")
			val, _ := reader.ReadString('\n')
			if strings.HasPrefix(val, "U") {
				s.up = CameraServoConfig{index: i, sign: 1.0}
				s.down = CameraServoConfig{index: i, sign: -1.0}
				break
			} else if strings.HasPrefix(val, "D") {
				s.up = CameraServoConfig{index: i, sign: -1.0}
				s.down = CameraServoConfig{index: i, sign: 1.0}
				break
			} else if strings.HasPrefix(val, "L") {
				s.left = CameraServoConfig{index: i, sign: 1.0}
				s.right = CameraServoConfig{index: i, sign: -1.0}
				break
			} else if strings.HasPrefix(val, "R") {
				s.left = CameraServoConfig{index: i, sign: -1.0}
				s.right = CameraServoConfig{index: i, sign: 1.0}
				break
			}
		}
		s.servos[i] = *servo
	}

	cfg, section, err := doLoadIniWithMatchingSectionOrCreateEmptyForCamera()
	if err != nil {
		return fmt.Errorf("Could not Load/Save Ini file. Error: %v", err)
	}

	section.NewKey("CameraUpIndex", fmt.Sprintf("%v", s.up.index))
	section.NewKey("CameraUpSign", fmt.Sprintf("%v", s.up.sign))

	section.NewKey("CameraDownIndex", fmt.Sprintf("%v", s.down.index))
	section.NewKey("CameraDownSign", fmt.Sprintf("%v", s.down.sign))

	section.NewKey("CameraLeftIndex", fmt.Sprintf("%v", s.left.index))
	section.NewKey("CameraLeftSign", fmt.Sprintf("%v", s.left.sign))

	section.NewKey("CameraRightIndex", fmt.Sprintf("%v", s.right.index))
	section.NewKey("CameraRightSign", fmt.Sprintf("%v", s.right.sign))

	err = cfg.SaveTo(CalibrationFilePath)
	if err != nil {
		return fmt.Errorf("Could not store configuration for camera. Error: %v", err)
	}

	return nil
}

//MoveUp will move the camera in an up position
func (s *CalibratedCamera) MoveUp(percent float64) error {
	return s.servos[s.up.index].Move(percent, s.up.sign)
}

//MoveDown will move the camera in a down position
func (s *CalibratedCamera) MoveDown(percent float64) error {
	return s.servos[s.down.index].Move(percent, s.down.sign)
}

//MoveLeft will move the camera in a left position
func (s *CalibratedCamera) MoveLeft(percent float64) error {
	return s.servos[s.left.index].Move(percent, s.left.sign)
}

//MoveRight will move the camera in a right position
func (s *CalibratedCamera) MoveRight(percent float64) error {
	return s.servos[s.right.index].Move(percent, s.right.sign)
}
