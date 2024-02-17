package gphoto2cli

import "sync"

// Port describe a specifi camera detected by gphoto2
type Port struct {
	Model string
	Port  string
	sync.Mutex
}

// CameraSummary detailed information about a camera
type CameraSummary struct {
	Manufacturer   string
	Model          string
	Version        string
	SerialNumber   string
	CaptureFormats []string
	DisplayFormats []string
}

// ConfigValue config value from a camera
type ConfigValue struct {
	Label    string
	Readonly bool
	Type     string
	Current  string
	Choices  map[int]string
}

// CaptureResult is the result of any capture action.
type CaptureResult struct {
	// OnCamera filename on the camera, only supported by Camera.CapturePhoto
	OnCamera string
	// LocalFile filename of the capture result
	LocalFile string
	// Delete if the file was deleted from the camera, only supported by Camera.CapturePhoto
	Deleted bool
}
