package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/blackjack/webcam"

	"sdmimaye.de/smart-video-car/car"
	"sdmimaye.de/smart-video-car/hardware"
	"sdmimaye.de/smart-video-car/stream"
)

func main() {
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()

	log.Println("Starting new Smart-Video-Car instance...")
	err = hardware.InitializeServoController()
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
		log.Panicf("Could not create new smart car instance. Error: %v", err)
	}
	execution := flag.String("e", "console", "The execution type of the application. Valid values are: console or tcp")
	flag.Parse()

	var s stream.Stream
	if execution == nil || strings.HasPrefix(*execution, "console") { //fallback to console
		log.Println("Will start console execution...")
		s = stream.ConsoleStream{}
	} else if strings.HasPrefix(*execution, "tcp") {
		log.Println("Will start tcp execution...")
		s, err = stream.NewTCPStream(1337)
		if err != nil {
			log.Panicf("Could not start new TCP Server on port 1337")
		}
	} else {
		log.Panicf("Unknown execution type: %v. Will exit now! (Valid values are: console or tcp)\n", *execution)
	}

	log.Printf("Exeuction: %v\n", *execution)
	car.Listen(s)
}
