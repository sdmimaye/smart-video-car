package steering

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
)

//HMovement represents a horizontal movement for the car
type HMovement int8

const (
	//HMovementNone represents movement in no particular direction
	HMovementNone HMovement = 0
	//HMovementLeft represents movement in a left direction
	HMovementLeft HMovement = 1
	//HMovementRight represents movement in a right direction
	HMovementRight HMovement = 2
)

//VMovement represents a vertical movement for the car
type VMovement int8

const (
	//VMovementNone represents movement in no particular direction
	VMovementNone VMovement = 0
	//VMovementUp represents movement in a up direction
	VMovementUp VMovement = 1
	//VMovementDown represents movement in a down direction
	VMovementDown VMovement = 2
)

//Step represents a movement step with a fixed speed, a direction and a camera movement
type Step struct {
	Speed                 float64
	CarMovement           HMovement
	CarMovementPercentage float64
	CameraHMovement       HMovement
	CameraHPercentage     float64
	CameraVMovement       VMovement
	CameraVPercentage     float64
}

//ParseStep will parse a step or return an error
func ParseStep(command []byte) (*Step, error) {
	reader := bytes.NewReader(command)
	order := binary.BigEndian

	speed := float64(0)
	err := binary.Read(reader, order, &speed)
	if err != nil {
		return nil, errors.New("Could not read speed from command bytes. Error: " + err.Error())
	}

	direction := int8(0)
	err = binary.Read(reader, order, &direction)
	if err != nil {
		return nil, errors.New("Could not read direction from command bytes")
	}

	dirpercent := float64(0)
	err = binary.Read(reader, order, &dirpercent)
	if err != nil {
		return nil, errors.New("Could not read direction-percent from command bytes")
	}

	camupdown := int8(0)
	err = binary.Read(reader, order, &camupdown)
	if err != nil {
		return nil, errors.New("Could not read camera up/down from command bytes")
	}

	cudpercent := float64(0)
	err = binary.Read(reader, order, &cudpercent)
	if err != nil {
		return nil, errors.New("Could not read camera up/down-percent from command bytes")
	}

	camleftright := int8(0)
	err = binary.Read(reader, order, &camleftright)
	if err != nil {
		return nil, errors.New("Could not read camera left/right from command bytes")
	}

	clrpercent := float64(0)
	err = binary.Read(reader, order, &clrpercent)
	if err != nil {
		return nil, errors.New("Could not read camera left/right-percent from command bytes")
	}
	log.Printf("Speed: %v, Direction: %v, DirectionPerc: %v, CamUpDown: %v, CamUpDownPerc: %v, CamLeftRight: %v, CamLeftRightPerc: %v", speed, direction, dirpercent, camupdown, cudpercent, camleftright, clrpercent)

	return &Step{
		Speed:                 speed,
		CarMovement:           HMovement(direction),
		CarMovementPercentage: dirpercent,
		CameraHMovement:       HMovement(camleftright),
		CameraHPercentage:     clrpercent,
		CameraVMovement:       VMovement(camupdown),
		CameraVPercentage:     cudpercent,
	}, nil
}
