package steering

import (
	"fmt"
	"log"
	"net"
)

//UDPEngine will steer the car over an udp socket
type UDPEngine struct {
	socket *net.UDPConn
}

//StartEngine will start the UDP socket and wait for incomming commands
func (s UDPEngine) StartEngine(sc StepCallback) error {
	if s.socket != nil {
		log.Println("UDP-Steering Engine is already running...")
		return nil
	}

	addr := net.UDPAddr{Port: 1338, IP: net.ParseIP("0.0.0.0")}
	socket, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return fmt.Errorf("Could not start UDP-Steering Engine. Error: %v", err)
	}

	log.Printf("Listening for incomming UDP instructions on Port: %v\n", addr.Port)
	go func() {
		var buffer [64]byte
		for socket != nil {
			cnt, remote, err := socket.ReadFromUDP(buffer[:])
			if err == nil {
				log.Printf("UDP Message received: Sender: %v, Length: %v, Data: %v", remote, cnt, buffer[:cnt])
				step, err := ParseStep(buffer[:cnt])
				if err != nil {
					log.Printf("Error while handling command: %v\n", err)
				} else {
					sc(step)
				}
			} else {
				log.Printf("Error while receving UDP Commands: %v\n", err)
				s.EndEngine()
			}
		}
	}()

	return nil
}

//EndEngine will terminate the socket
func (s UDPEngine) EndEngine() error {
	if s.socket == nil {
		log.Println("UDP Steering Engine has ended already...")
		return nil
	}

	err := s.socket.Close()
	s.socket = nil
	if err != nil {
		return fmt.Errorf("Error while closing udp steering socket: %v", err)
	}

	return nil
}
