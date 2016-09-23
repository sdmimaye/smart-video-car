package car

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"sdmimaye.de/smart-video-car/stream"
)

func doCalibrate(c *Car, stream stream.Stream) error {
	err := c.Motor.Calibrate(stream)
	if err != nil {
		return err
	}

	err = c.Steering.Calibrate(stream)
	if err != nil {
		return err
	}

	err = c.Camera.Calibrate(stream)
	if err != nil {
		return err
	}

	return nil
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
	fmt.Fprint(w, "Press any key to exit steering\r\n")
	reader.ReadString('\n')
	err = s.EndEngine()
	if err != nil {
		return fmt.Errorf("Could not stop steering method: %v. Error: %v", command, err)
	}
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
			}
		}
	})
}
