package car

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"sdmimaye.de/smart-video-car/steering"
	"sdmimaye.de/smart-video-car/stream"
)

type calibration struct {
	motor    bool
	steering bool
	camera   bool
}

func doCalibrate(c *Car, stream stream.Stream) error {
	r := stream.GetReader()
	w := stream.GetWriter()

	for {
		_, err := fmt.Fprint(w, "Please choose which part you want to calibrate:\r\n[0] Everything\r\n[1] Motor\r\n[2] Steering\r\n[3] Camera\r\nAnything else will bring you back to the previous selection\r\n")
		reader := bufio.NewReader(r)
		command, _ := reader.ReadString('\n')
		var cali calibration

		if strings.HasPrefix(command, "0") {
			cali = calibration{true, true, true}
		} else if strings.HasPrefix(command, "1") {
			cali = calibration{true, false, false}
		} else if strings.HasPrefix(command, "2") {
			cali = calibration{false, true, false}
		} else if strings.HasPrefix(command, "3") {
			cali = calibration{false, false, true}
		}

		if cali.motor {
			err = c.Motor.Calibrate(stream)
			if err != nil {
				return err
			}
		}
		if cali.steering {
			err = c.Steering.Calibrate(stream)
			if err != nil {
				return err
			}
		}
		if cali.camera {
			err = c.Camera.Calibrate(stream)
			if err != nil {
				return err
			}
		}
		if !cali.camera && !cali.motor && !cali.steering {
			return nil
		}
	}
}

func doSteer(c *Car, stream stream.Stream) error {
	r := stream.GetReader()
	w := stream.GetWriter()
	reader := bufio.NewReader(r)

	var s steering.Engine
	fmt.Fprint(w, "Please select your steering method:\r\n[0] UDP\r\n")
	command, _ := reader.ReadString('\n')
	if strings.HasPrefix(command, "0") {
		s = steering.UDPEngine{}
	} else {
		return fmt.Errorf("Unknown steering method: %v\r\n", command)
	}

	err := s.StartEngine(func(step *steering.Step) {
		err := c.Move(step)
		if err != nil {
			log.Printf("Could not steer car. Error: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("Could not start steering method: %v. Error: %v", command, err)
	}
	defer s.EndEngine()
	fmt.Fprint(w, "Press any key to exit steering\r\n")
	reader.ReadString('\n')
	return nil
}

//Execute will execute all car related commands
func Execute(c *Car, stream stream.Stream) {
	stream.OnConnectionEstablished(func() {
		log.Println("(Re-)Starting Car-Execution")
		r := stream.GetReader()
		w := stream.GetWriter()

		for {
			fmt.Fprint(w, "Please enter your next command:\r\n[0] Calibrate\r\n[1] Steer\r\n")
			reader := bufio.NewReader(r)

			command, _ := reader.ReadString('\n')
			if strings.HasPrefix(command, "0") {
				err := doCalibrate(c, stream)
				if err != nil {
					fmt.Fprintf(w, "Error while calibrating car. Error: %v\r\n", err)
				}
			} else if strings.HasPrefix(command, "1") {
				err := doSteer(c, stream)
				if err != nil {
					fmt.Fprintf(w, "Error while steering car. Error: %v\r\n", err)
				}
			} else {
				return
			}
		}
	})
}
