package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"sdmimaye.de/smart-video-car/car"
	"sdmimaye.de/smart-video-car/hardware"
	"sdmimaye.de/smart-video-car/stream"
)

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

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		hardware.DeInitializeMotorController()
		hardware.DeInitializeServoController()
		os.Exit(1)
	}()

	car, err := car.NewCar()
	if err != nil {
		log.Panicf("Cloud not create new smart car instance. Error: %v", err)
	}
	/*
		s, err := stream.NewTCPStream(1337)
		if err != nil {
			log.Panicf("Could not start new TCP Server on port 1337")
		}
		car.Listen(*s)
	*/
	car.Listen(stream.ConsoleStream{})
}
