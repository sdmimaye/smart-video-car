package car

import (
	"errors"
	"fmt"

	"sdmimaye.de/smart-video-car/components"
	"sdmimaye.de/smart-video-car/steering"
	"sdmimaye.de/smart-video-car/stream"
)

//Car represents our smart-video-car
type Car struct {
	Motor    *components.CalibratedMotor
	Steering *components.CalibratedSteering
	Camera   *components.CalibratedCamera
}

//NewCar will create a new smart car instance
func NewCar() (*Car, error) {
	motor, err := components.NewCalibratedMotor()
	if err != nil {
		return nil, errors.New("Could not create calibrated Motor for car. Error: " + err.Error())
	}
	steering, err := components.NewCalibratedSteering()
	if err != nil {
		return nil, errors.New("Could not create calibrated Steering for car. Error: " + err.Error())
	}

	camera, err := components.NewCalibratedCamera()
	if err != nil {
		return nil, errors.New("Could not create calibrated Camera for car. Error: " + err.Error())
	}

	return &Car{Motor: motor, Camera: camera, Steering: steering}, nil
}

//Listen will make the car listen to the incomming requests from the stream and move accordingly
func (c *Car) Listen(stream stream.Stream) {
	Execute(c, stream)
}

//Move will move the car with a certain step
func (c *Car) Move(step *steering.Step) error {
	err := c.Motor.SetSpeed(step.Speed)
	if err != nil {
		return fmt.Errorf("Could not accelerate/decelerate. Error: %v", err)
	}

	switch step.CarMovement {
	case steering.HMovementNone:
		err := c.Steering.Center()
		if err != nil {
			return fmt.Errorf("Could not center steer. Error: %v", err)
		}
	case steering.HMovementLeft:
		err = c.Steering.SteerLeft(step.CarMovementPercentage)
		if err != nil {
			return fmt.Errorf("Could not steer left. Error: %v", err)
		}
	case steering.HMovementRight:
		err = c.Steering.SteerRight(step.CarMovementPercentage)
		if err != nil {
			return fmt.Errorf("Could not steer right. Error: %v", err)
		}
	default:
		return fmt.Errorf("Unknown direction: %v. Use either none(0), left(1) or right(2)", step.CarMovement)
	}

	switch step.CameraVMovement {
	case steering.VMovementNone:
		err := c.Camera.CenterUpDown()
		if err != nil {
			return fmt.Errorf("Could not center camera up/down. Error: %v", err)
		}
	case steering.VMovementUp:
		err = c.Camera.MoveUp(step.CameraVPercentage)
		if err != nil {
			return fmt.Errorf("Could not move camera up. Error: %v", err)
		}
	case steering.VMovementDown:
		err = c.Camera.MoveDown(step.CameraVPercentage)
		if err != nil {
			return fmt.Errorf("Could not move camera down. Error: %v", err)
		}
	default:
		return fmt.Errorf("Unknown camera up/down direction: %v. Use either none(0), up(1) or down(2)", step.CameraVMovement)
	}

	switch step.CameraHMovement {
	case steering.HMovementNone:
		err := c.Camera.CenterLeftRight()
		if err != nil {
			return fmt.Errorf("Could not center camera left/right. Error: %v", err)
		}
	case steering.HMovementLeft:
		err = c.Camera.MoveLeft(step.CameraHPercentage)
		if err != nil {
			return fmt.Errorf("Could not move camera left. Error: %v", err)
		}
	case steering.HMovementRight:
		err = c.Camera.MoveRight(step.CameraHPercentage)
		if err != nil {
			return fmt.Errorf("Could not move camera right. Error: %v", err)
		}
	default:
		return fmt.Errorf("Unknown cameraleft/right direction: %v. Use either left(0) or right(1)", step.CameraHMovement)
	}

	return nil
}
