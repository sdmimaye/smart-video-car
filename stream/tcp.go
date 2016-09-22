package stream

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var listener net.Listener
var connection net.Conn
var reset func()

//TCPStream represents a Tcp-Server which will listen on a specified port to communicate with this system
type TCPStream struct {
}

func (t TCPStream) Read(p []byte) (int, error) {
	for connection == nil {
		time.Sleep(1 + time.Second)
	}

	n, err := connection.Read(p)
	if err != nil {
		t.Close()
		return 0, nil
	}

	return n, nil
}

func (t TCPStream) Write(p []byte) (int, error) {
	for connection == nil {
		time.Sleep(1 + time.Second)
	}

	n, err := connection.Write(p)
	if err != nil {
		t.Close()
		return 0, nil
	}

	return n, nil
}

//NewTCPStream will create a new TCPStream or create an error
func NewTCPStream(port int) (*TCPStream, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("Could not generate TCPStream on Port: %v. Error: %v", port, err)
	}

	stream := TCPStream{}
	go func() {
		for {
			log.Print("Waiting for incomming TCP connection...\n")
			c, err := listener.Accept()
			if err != nil {
				log.Panicf("Error while accepting TCP connection. Error: %v\n", err)
			}
			go func(socket net.Conn) {
				stream.Close()
				connection = socket
			}(c)
		}
	}()
	return &stream, nil
}

//GetReader will return the reading stream
func (t TCPStream) GetReader() io.Reader {
	return t
}

//GetWriter will return the reading stream
func (t TCPStream) GetWriter() io.Writer {
	return t
}

//Close will will close the tcp stream
func (t TCPStream) Close() error {
	if connection == nil {
		return nil
	}

	connection.Close()
	connection = nil
	if reset != nil {
		reset()
	}

	return nil
}

//OnConnectionEstablished can be set to signalize a reset in the communication
func (t TCPStream) OnConnectionEstablished(f func()) {
	reset = f
	f()
}
