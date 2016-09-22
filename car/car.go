package car

import (
	"errors"

	"sdmimaye.de/smart-video-car/components"
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
