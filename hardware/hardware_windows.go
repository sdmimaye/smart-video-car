package hardware

import (
	"errors"
	"log"
)

var servoInitialized = false
var motorInitialized = false

//WindowsServoMotor represents a Servo-Motor inside a linux enviroment
type WindowsServoMotor struct {
	channel int
}

//WindowsPin represents a (fake Windows) digital GPIO Pin
type WindowsPin struct {
	pin int
}

//SetAngle will set the angle of a windows servo motor
func (s WindowsServoMotor) SetAngle(angle int) error {
	log.Printf("Setting Angle of Servo (with Channel:%v) to: %v\n", s.channel, angle)
	return nil
}

//SetDirection will set the diretion (of this fake device) to input or output
func (p WindowsPin) SetDirection(direction PinDirection) error {
	log.Printf("Setting direction of pin: %v, to: %v\n", p.pin, direction)
	return nil
}

//Read will read the current level of a (fake Windows) GPIO pin
func (p WindowsPin) Read() (PinLevel, error) {
	log.Printf("Reading value of pin: %v\n", p.pin)
	return Low, nil
}

//Write will read the current level of a GPIO pin
func (p WindowsPin) Write(level PinLevel) error {
	log.Printf("Writing value: %v of pin: %v\n", level, p.pin)
	return nil
}

//InitializeServoController will initialize the (fake Windows) I²C Bus
func InitializeServoController() error {
	log.Print("Initializing Windows Servo Controller...\n")
	servoInitialized = true
	return nil
}

//InitializeMotorController will initialize the GPIO-Pins responsible for the motor movement
func InitializeMotorController() error {
	log.Print("Initializing Windows Motor Controller...\n")
	motorInitialized = true
	return nil
}

//DeInitializeServoController will deinitialize the (fake Windows) I²C bus
func DeInitializeServoController() error {
	log.Print("Deinitializing Windows Servo Controller...\n")
	servoInitialized = false
	return nil
}

//DeInitializeMotorController will deinitialize the (fake Windows) GPIO-Pins responsible for the motor movement
func DeInitializeMotorController() error {
	log.Print("Deinitializing Windows Motor Controller...\n")
	motorInitialized = false
	return nil
}

//GetServo will create a servo out of a channel
func GetServo(channel int) (ServoMotor, error) {
	if !servoInitialized {
		return nil, errors.New("ServoController is not initialized...")
	}

	log.Printf("Generated new Servo-Motor on channel: %v\n", channel)
	return WindowsServoMotor{channel: channel}, nil
}

//GetPin will create a new GPIO Pin
func GetPin(pin int) (Pin, error) {
	log.Printf("Generating GPIO pin: %v\n", pin)
	return WindowsPin{pin: pin}, nil
}

//SetPwmValue will set the PWM Value on a channel
func SetPwmValue(channel int, onTime int, offTime int) error {
	if !servoInitialized {
		return errors.New("Please initialize the (fake Windows) I²C Controller before using pwm")
	}

	log.Printf("Setting (fake Windows) PWM Signal. Channel: %v, on: %v, off: %v\n", channel, onTime, offTime)
	return nil
}
