package hardware

//ServoMotor will abstract a servor motor controlled by an servo controller in this application
type ServoMotor interface {
	SetAngle(angle int) error
}

//PinDirection determines the direction for the communication for one GPIO pin
type PinDirection int

//PinLevel is the level high or low on a pin
type PinLevel int

const (
	// In represents read mode.
	In PinDirection = iota
	// Out represents write mode.
	Out
)

const (
	// Low represents 0.
	Low PinLevel = iota

	// High represents 1.
	High
)

//Pin represents a GPIO Pin inside our device
type Pin interface {
	SetDirection(direction PinDirection) error
	Read() (PinLevel, error)
	Write(level PinLevel) error
}
