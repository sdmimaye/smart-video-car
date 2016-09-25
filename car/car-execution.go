package car

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"sdmimaye.de/smart-video-car/stream"
)

type execution struct {
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
		var exe execution

		if strings.HasPrefix(command, "0") {
			exe = execution{true, true, true}
		} else if strings.HasPrefix(command, "1") {
			exe = execution{true, false, false}
		} else if strings.HasPrefix(command, "2") {
			exe = execution{false, true, false}
		} else if strings.HasPrefix(command, "3") {
			exe = execution{false, false, true}
		}

		if exe.motor {
			err = c.Motor.Calibrate(stream)
			if err != nil {
				return err
			}
		}
		if exe.steering {
			err = c.Steering.Calibrate(stream)
			if err != nil {
				return err
			}
		}
		if exe.camera {
			err = c.Camera.Calibrate(stream)
			if err != nil {
				return err
			}
		}
		if !exe.camera && !exe.motor && !exe.steering {
			return nil
		}
	}
}

func doSteer(c *Car, stream stream.Stream) error {
	r := stream.GetReader()
	w := stream.GetWriter()
	reader := bufio.NewReader(r)

	var s SteeringEngine
	fmt.Fprint(w, "Please select your steering method:\r\n[0] UDP\r\n")
	command, _ := reader.ReadString('\n')
	if strings.HasPrefix(command, "0") {
		s = UDPSteeringEngine{}
	} else {
		return fmt.Errorf("Unknown steering method: %v\r\n", command)
	}

	err := s.StartEngine(c)
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
