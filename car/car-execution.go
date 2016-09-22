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

			command, err := reader.ReadString('\n')
			if err != nil {
				stream.Close()
			}

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
