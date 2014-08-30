package devices

// DeviceRepository holds all the methods needed to save, delete, load devices.
type DeviceRepository interface {
	Load(user string) ([]*Device, error)
	Delete(hash string) error
	Save(device *Device) error
	List(users []string) ([]*Device, error)
	Push(users []string, message string) error
}
