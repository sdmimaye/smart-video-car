package main

import "log"
import "sdmimaye.de/smart-video-car/hardware"
import "sdmimaye.de/smart-video-car/car"
import "sdmimaye.de/smart-video-car/stream"

func main() {
	log.Println("Starting new Smart-Video-Car instance...")
	err := hardware.InitializeServoController()
	if err != nil {
		log.Panicf("Could not initialize Servo-Controller: %v\r\n", err)
	}
	defer hardware.DeInitializeServoController()

	err = hardware.InitializeMotorController()
	if err != nil {
		log.Panicf("Could not initialize Motor-Controller: %v\r\n", err)
	}
	defer hardware.DeInitializeMotorController()

	car, err := car.NewCar()
	if err != nil {
		log.Panicf("Cloud not create new smart car instance. Error: %v", err)
	}

	s, err := stream.NewTCPStream(1337)
	if err != nil {
		log.Panicf("Could not start new TCP Server on port 1337")
	}
	car.Listen(*s)
}
