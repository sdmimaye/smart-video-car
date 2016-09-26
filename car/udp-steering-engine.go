package car

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)

var socket *net.UDPConn

//UDPSteeringEngine will steer the car over an udp socket
type UDPSteeringEngine struct {
}

//StartEngine will start the UDP socket and wait for incomming commands
func (s UDPSteeringEngine) StartEngine(c *Car) error {
	if socket != nil {
		log.Println("UDp-Steering Engine is already running...")
		return nil
	}

	addr := net.UDPAddr{Port: 1338, IP: net.ParseIP("0.0.0.0")}
	socket, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return fmt.Errorf("Could not start UDP-Steering Engine. Error: %v", err)
	}

	log.Printf("Listening for incomming UDP instructions on Port: %v\n", addr.Port)
	go func() {
		var buffer [32]byte
		for socket != nil {
			cnt, remote, err := socket.ReadFromUDP(buffer[:])
			if err == nil {
				log.Printf("UDP Message received: Sender: %v, Length: %v, Data: %v", remote, cnt, buffer[:cnt])
				err = doHandleCommand(c, buffer[:cnt])
				if err != nil {
					log.Printf("Error while handling command: %v\n", err)
				}
				/*
					table := crc8.MakeTable(crc8.CRC8)
					chksum := crc8.Checksum(buffer[:31], table)
					if chksum == buffer[31] {
						doHandleCommand(c, buffer[:cnt])
					} else {
						log.Printf("Checksum Error, got: %v, expected: %v", buffer[31], chksum)
					}
				*/
			} else {
				log.Printf("Error while receving UDP Commands: %v\n", err)
				s.EndEngine()
			}
		}
	}()

	return nil
}

func doHandleCommand(c *Car, command []byte) error {
	reader := bytes.NewReader(command)
	order := binary.BigEndian

	speed := float64(0)
	err := binary.Read(reader, order, &speed)
	if err != nil {
		return errors.New("Could not read speed from command bytes. Error: " + err.Error())
	}

	direction := int8(0)
	err = binary.Read(reader, order, &direction)
	if err != nil {
		return errors.New("Could not read direction from command bytes")
	}

	dirpercent := float64(0)
	err = binary.Read(reader, order, &dirpercent)
	if err != nil {
		return errors.New("Could not read direction-percent from command bytes")
	}

	camupdown := int8(0)
	err = binary.Read(reader, order, &camupdown)
	if err != nil {
		return errors.New("Could not read camera up/down from command bytes")
	}

	cudpercent := float64(0)
	err = binary.Read(reader, order, &cudpercent)
	if err != nil {
		return errors.New("Could not read camera up/down-percent from command bytes")
	}

	camleftright := int8(0)
	err = binary.Read(reader, order, &camleftright)
	if err != nil {
		return errors.New("Could not read camera left/right from command bytes")
	}

	clrpercent := float64(0)
	err = binary.Read(reader, order, &clrpercent)
	if err != nil {
		return errors.New("Could not read camera left/right-percent from command bytes")
	}

	log.Printf("Speed: %v, Direction: %v, CamUpDown: %v, CamUpDownPerc: %v, CamLeftRight: %v, CamLeftRightPerc: %v", speed, direction, camupdown, cudpercent, camleftright, clrpercent)

	err = c.Motor.SetSpeed(speed)
	if err != nil {
		return fmt.Errorf("Could not accelerate/decelerate. Error: %v", err)
	}

	switch direction {
	case 0: //nothing
		err := c.Steering.Center()
		if err != nil {
			return fmt.Errorf("Could not center steer. Error: %v", err)
		}
	case 1: //left
		err = c.Steering.SteerLeft(dirpercent)
		if err != nil {
			return fmt.Errorf("Could not steer left. Error: %v", err)
		}
	case 2: //right
		err = c.Steering.SteerRight(dirpercent)
		if err != nil {
			return fmt.Errorf("Could not steer right. Error: %v", err)
		}
	default:
		return fmt.Errorf("Unknown direction: %v. Use either left(0) or right(1)", direction)
	}

	switch camupdown {
	case 0: //home
		err := c.Camera.CenterUpDown()
		if err != nil {
			return fmt.Errorf("Could not center camera up/down. Error: %v", err)
		}
	case 1: //up
		err = c.Camera.MoveUp(cudpercent)
		if err != nil {
			return fmt.Errorf("Could not move camera up. Error: %v", err)
		}
	case 2: //down
		err = c.Camera.MoveDown(cudpercent)
		if err != nil {
			return fmt.Errorf("Could not move camera down. Error: %v", err)
		}
	default:
		return fmt.Errorf("Unknown camera up/down direction: %v. Use either up(0) or down(1)", camupdown)
	}

	switch camleftright {
	case 0: //home
		err := c.Camera.CenterLeftRight()
		if err != nil {
			return fmt.Errorf("Could not center camera left/right. Error: %v", err)
		}
	case 1: //left
		err = c.Camera.MoveLeft(clrpercent)
		if err != nil {
			return fmt.Errorf("Could not move camera left. Error: %v", err)
		}
	case 2: //right
		err = c.Camera.MoveRight(clrpercent)
		if err != nil {
			return fmt.Errorf("Could not move camera right. Error: %v", err)
		}
	default:
		return fmt.Errorf("Unknown cameraleft/right direction: %v. Use either left(0) or right(1)", camleftright)
	}

	return nil
}

//EndEngine will terminate the socket
func (s UDPSteeringEngine) EndEngine() error {
	if socket == nil {
		log.Println("UDP Steering Engine has ended already...")
		return nil
	}

	err := socket.Close()
	socket = nil
	if err != nil {
		return fmt.Errorf("Error while closing udp steering socket: %v", err)
	}

	return nil
}
