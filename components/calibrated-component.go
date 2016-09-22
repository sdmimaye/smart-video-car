package components

import "sdmimaye.de/smart-video-car/stream"

//CalibratedComponent represents a car component which needs calibration
type CalibratedComponent interface {
	Calibrate(stream stream.Stream) error
}
