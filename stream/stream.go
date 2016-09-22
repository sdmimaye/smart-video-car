package stream

import "io"

//Stream is used to decouple the input/output mechanisms
type Stream interface {
	GetReader() io.Reader
	GetWriter() io.Writer
	OnConnectionEstablished(f func())
	Close() error
}
