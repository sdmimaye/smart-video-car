package hardware

import (
	"errors"
	"fmt"
	"log"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"

	_ "github.com/kidoman/embd/host/all" //Otherwise this packae will not be loaded
	"github.com/kidoman/embd/motion/servo"
)

var pca *pca9685.PCA9685

//LinuxServoMotor represents a Servo-Motor inside a linux enviroment
type LinuxServoMotor struct {
	servo *servo.Servo
}

//LinuxPin represents a digital GPIO Pin
type LinuxPin struct {
	pin embd.DigitalPin
}

//SetAngle will set the angle of a linux servo motor
func (s LinuxServoMotor) SetAngle(angle int) error {
	return s.servo.SetAngle(angle)
}

//SetDirection will set the diretion to input or output
func (p LinuxPin) SetDirection(direction PinDirection) error {
	dir := embd.Direction(direction)
	return p.pin.SetDirection(dir)
}

//Read will read the current level of a GPIO pin
func (p LinuxPin) Read() (PinLevel, error) {
	level, err := p.pin.Read()
	if err != nil {
		return Low, fmt.Errorf("Could not read from GPIO pin. Error: %v", err)
	}

	return PinLevel(level), nil
}

//Write will read the current level of a GPIO pin
func (p LinuxPin) Write(level PinLevel) error {
	log.Printf("Writing: %v on gpio pin: %v", level, p.pin.N())
	err := p.pin.Write(int(level))
	if err != nil {
		return fmt.Errorf("Could not write to GPIO pin. Error: %v", err)
	}

	return nil
}

//InitializeServoController will initialize the I²C bus
func InitializeServoController() error {
	if pca != nil {
		log.Println("I²C is already initialized... skipped")
		return nil
	}

	log.Println("Initializing I²C Bus")
	err := embd.InitI2C()
	if err != nil {
		return fmt.Errorf("Could not initialize I²C Bus. Reason: %v", err)
	}

	bus := embd.NewI2CBus(1)

	pca = pca9685.New(bus, 0x40)
	pca.Freq = 50
	err = pca.Wake()
	if err != nil {
		return fmt.Errorf("Could not wake pwm I²C Bus. Reason: %v", err)
	}

	return nil
}

//InitializeMotorController will initialize the GPIO-Pins responsible for the motor movement
func InitializeMotorController() error {
	log.Println("Initializing GPIO pins")
	err := embd.InitGPIO()
	if err != nil {
		return fmt.Errorf("Could not initialize motor controller. Error: %v", err)
	}

	return nil
}

//DeInitializeServoController will deinitialize the I²C bus
func DeInitializeServoController() error {
	log.Println("Closing I²C...")
	if pca == nil {
		log.Println("I²C is already deinitialized... skipped")
		return nil
	}

	err := pca.Sleep()
	if err != nil {
		return fmt.Errorf("Could set sleep PCA. Error: %v", err)
	}

	err = pca.Close()
	if err != nil {
		return fmt.Errorf("Could not close PCA. Error: %v", err)
	}

	err = embd.CloseI2C()
	if err != nil {
		return fmt.Errorf("Could not deinitialize I²C Bus. Error: %v", err)
	}

	return nil
}

//DeInitializeMotorController will deinitialize the GPIO-Pins responsible for the motor movement
func DeInitializeMotorController() error {
	log.Println("Closing GPIO...")
	err := embd.CloseGPIO()
	if err != nil {
		return fmt.Errorf("Could not deinitialize motor controller. Error: %v", err)
	}

	return nil
}

//GetServo will create a servo out of a channel
func GetServo(channel int) (ServoMotor, error) {
	log.Printf("Generating I²C Servo on channel: %v\n", channel)
	if pca == nil {
		return nil, errors.New("PCA is not initialized")
	}

	ch := pca.ServoChannel(channel)
	servo := servo.New(ch)

	return LinuxServoMotor{servo: servo}, nil
}

//GetPin will create a new GPIO Pin
func GetPin(pin int) (Pin, error) {
	log.Printf("Generating GPIO pin: %v\n", pin)
	dpin, err := embd.NewDigitalPin(pin)
	if err != nil {
		return nil, fmt.Errorf("Error while opening GPIO pin: %v. Error: %v", pin, err)
	}

	return LinuxPin{pin: dpin}, nil
}

//SetPwmValue will set the PWM Value on a channel
func SetPwmValue(channel int, onTime int, offTime int) error {
	if pca == nil {
		return errors.New("Please initialize the I²C Controller before using pwm")
	}

	log.Printf("Setting PWM on channel: %v, on: %v, off: %v\n", channel, onTime, offTime)
	return pca.SetPwm(channel, onTime, offTime)
}
