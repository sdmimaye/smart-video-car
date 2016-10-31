package steering

//StepCallback is a callback function which will be called when a step occured
type StepCallback func(*Step)

//Engine will start a new steering implementation
type Engine interface {
	StartEngine(sc StepCallback) error
	EndEngine() error
}
