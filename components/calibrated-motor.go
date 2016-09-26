package components

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/go-ini/ini"

	"sdmimaye.de/smart-video-car/hardware"
	"sdmimaye.de/smart-video-car/stream"
)

const (
	motor0Pin0     = 17
	motor0Pin1     = 18
	motor1Pin0     = 27
	motor1Pin1     = 22
	speedPwmMotor1 = 4
	speedPwmMotor0 = 5
)

type motorCabling int

const (
	p0ForwardP1Backward motorCabling = 0
	p1ForwardP0Backward motorCabling = 1
)

type motor struct {
	p0              hardware.Pin
	p1              hardware.Pin
	cabling         motorCabling
	speedPwmChannel int
}

//CalibratedMotor represents a calibrated motor inside our smart car. It composes out of 4 gpio pins and 2 pwn signals
type CalibratedMotor struct {
	m0 motor
	m1 motor
}

func doLoadIniWithMatchingSectionOrCreateEmptyForMotor() (*ini.File, *ini.Section, error) {
	cfg, err := ini.LooseLoad(CalibrationFilePath)
	if err != nil {
		return nil, nil, errors.New("Could not Load ini file. Error: " + err.Error())
	}

	sec, err := cfg.GetSection("Motor")
	if err != nil {
		log.Print("Missing Section: Motor. Will regenerate section\n")
		sec, err = cfg.NewSection("Motor")
		if err != nil {
			return nil, nil, errors.New("Error while generating Section. Error: " + err.Error())
		}
	}

	if !sec.HasKey("M0Cabling") {
		sec.NewKey("M0Cabling", "0")
	}

	if !sec.HasKey("M1Cabling") {
		sec.NewKey("M1Cabling", "0")
	}

	cfg.SaveTo(CalibrationFilePath)

	return cfg, sec, nil
}

func getPinAndSetToOutput(pin int) (hardware.Pin, error) {
	p, err := hardware.GetPin(pin)
	if err != nil {
		return nil, fmt.Errorf("Could not create GPIO Pin: %v for Motor", pin)
	}
	err = p.SetDirection(hardware.Out)
	if err != nil {
		log.Printf("Warning. Could not set GPIO direction for pin: %v. If the application was started a second time everything is fine. Otherwise a hardware error occured. Error: %v", pin, err)
	}

	return p, nil
}

//NewCalibratedMotor will create a new calibrated motor or return an error
func NewCalibratedMotor() (*CalibratedMotor, error) {
	_, section, err := doLoadIniWithMatchingSectionOrCreateEmptyForMotor()
	if err != nil {
		return nil, fmt.Errorf("Could not load motor config from ini. Error: %v", err)
	}

	m0p0, err := getPinAndSetToOutput(motor0Pin0)
	if err != nil {
		return nil, err
	}

	m0p1, err := getPinAndSetToOutput(motor0Pin1)
	if err != nil {
		return nil, err
	}
	m1p0, err := getPinAndSetToOutput(motor1Pin0)
	if err != nil {
		return nil, err
	}
	m1p1, err := getPinAndSetToOutput(motor1Pin1)
	if err != nil {
		return nil, err
	}

	m0cab, _ := section.Key("M0Cabling").Int()
	m1cab, _ := section.Key("M1Cabling").Int()

	m0 := motor{p0: m0p0, p1: m0p1, speedPwmChannel: speedPwmMotor0, cabling: motorCabling(m0cab)}
	m1 := motor{p0: m1p0, p1: m1p1, speedPwmChannel: speedPwmMotor1, cabling: motorCabling(m1cab)}
	motor := CalibratedMotor{m0: m0, m1: m1}
	return &motor, nil
}

//Calibrate will (re)calibrate a motor
func (m *CalibratedMotor) Calibrate(stream stream.Stream) error {
	r := stream.GetReader()
	w := stream.GetWriter()

	reader := bufio.NewReader(r)
	fmt.Fprint(w, "The first wheel will move in one direction. Please pay attention!\r\n")
	err := hardware.SetPwmValue(m.m0.speedPwmChannel, 0, 2000)
	if err != nil {
		return fmt.Errorf("Could not set speed via pwm channel: %v. Error: %v", m.m0.speedPwmChannel, err)
	}
	err = hardware.SetPwmValue(m.m1.speedPwmChannel, 0, 2000)
	if err != nil {
		return fmt.Errorf("Could not set speed via pwm channel: %v. Error: %v", m.m1.speedPwmChannel, err)
	}

	err = m.m0.p0.Write(hardware.High)
	if err != nil {
		return fmt.Errorf("Could not set pin %v to high. Error: %v", motor0Pin0, err)
	}
	fmt.Fprint(w, "In which direction is the wheel moving?\r\n[0] Forward\r\n[1] Backward\r\n")
	val, _ := reader.ReadString('\n')

	if strings.HasPrefix(val, "0") {
		m.m0.cabling = p0ForwardP1Backward
	} else if strings.HasPrefix(val, "1") {
		m.m0.cabling = p1ForwardP0Backward
	} else {
		fmt.Fprint(w, "Skipping configuration for first wheel...\r\n")
	}
	err = m.m0.p0.Write(hardware.Low)
	if err != nil {
		return fmt.Errorf("Could not set pin %v to low. Error: %v", motor0Pin0, err)
	}

	fmt.Fprint(w, "The second wheel will move in one direction. Please pay attention!\r\n")
	err = m.m1.p0.Write(hardware.High)
	if err != nil {
		return fmt.Errorf("Could not set pin %v to high. Error: %v", motor1Pin0, err)
	}
	fmt.Fprint(w, "In which direction is the wheel moving?\r\n[0] Forward\r\n[1] Backward\r\n")
	val, _ = reader.ReadString('\n')

	if strings.HasPrefix(val, "0") {
		m.m1.cabling = p0ForwardP1Backward
	} else if strings.HasPrefix(val, "1") {
		m.m1.cabling = p1ForwardP0Backward
	} else {
		fmt.Fprint(w, "Skipping configuration for second wheel...\r\n")
	}
	err = m.m1.p0.Write(hardware.Low)
	if err != nil {
		return fmt.Errorf("Could not set pin %v to low. Error: %v", motor1Pin0, err)
	}

	return nil
}

//SetSpeed will set the speed of the motor
func (m *CalibratedMotor) SetSpeed(speedPercentage float64) error {
	if speedPercentage > 100 || speedPercentage < -100 {
		return fmt.Errorf("Invalid speed percentage value: %v. Choose a valud between 100 and -100", speedPercentage)
	}

	pwm := int(math.Abs(40.96 * speedPercentage))
	err0 := hardware.SetPwmValue(m.m0.speedPwmChannel, 0, pwm)
	err1 := hardware.SetPwmValue(m.m1.speedPwmChannel, 0, pwm)

	if err0 != nil || err1 != nil {
		return fmt.Errorf("Could not set motor speed to: %v percent. Errors: %v, %v", speedPercentage, err0, err1)
	}

	if speedPercentage > 0 { //we are moving forward
		if m.m0.cabling == p0ForwardP1Backward {
			m.m0.p0.Write(hardware.High)
			m.m0.p1.Write(hardware.Low)
		} else if m.m0.cabling == p1ForwardP0Backward {
			m.m0.p0.Write(hardware.Low)
			m.m0.p1.Write(hardware.High)
		}

		if m.m1.cabling == p0ForwardP1Backward {
			m.m1.p0.Write(hardware.High)
			m.m1.p1.Write(hardware.Low)
		} else if m.m0.cabling == p1ForwardP0Backward {
			m.m1.p0.Write(hardware.Low)
			m.m1.p1.Write(hardware.High)
		}
	} else if speedPercentage < 0 { //we are moving backward
		if m.m0.cabling == p0ForwardP1Backward {
			m.m0.p0.Write(hardware.Low)
			m.m0.p1.Write(hardware.High)
		} else if m.m0.cabling == p1ForwardP0Backward {
			m.m0.p0.Write(hardware.High)
			m.m0.p1.Write(hardware.Low)
		}

		if m.m1.cabling == p0ForwardP1Backward {
			m.m1.p0.Write(hardware.Low)
			m.m1.p1.Write(hardware.High)
		} else if m.m0.cabling == p1ForwardP0Backward {
			m.m1.p0.Write(hardware.High)
			m.m1.p1.Write(hardware.Low)
		}
	} else { //we stopped
		m.m0.p0.Write(hardware.Low)
		m.m0.p1.Write(hardware.Low)
		m.m1.p0.Write(hardware.Low)
		m.m1.p1.Write(hardware.Low)
	}

	return nil
}

//Stop will halt all movemnt of the motor
func (m *CalibratedMotor) Stop() error {
	return m.SetSpeed(0)
}
