package car

//SteeringEngine will start a new steering implementation
type SteeringEngine interface {
	StartEngine(c *Car) error
	EndEngine() error
}
