package stream

import (
	"io"
	"os"
)

//ConsoleStream represents the console for generic input/output operations
type ConsoleStream struct {
}

//GetReader will return the reading stream
func (c ConsoleStream) GetReader() io.Reader {
	return os.Stdin
}

//GetWriter will return the reading stream
func (c ConsoleStream) GetWriter() io.Writer {
	return os.Stdout
}

//Close will will close the consoel stream
func (c ConsoleStream) Close() error {
	return nil
}

//OnConnectionEstablished can be set to signalize a reset in the communication
func (c ConsoleStream) OnConnectionEstablished(f func()) {
	f()
}
