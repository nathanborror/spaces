package rooms

// RoomRepository holds all the methods needed to save, delete, load and list User objects.
type RoomRepository interface {
	Load(hash string) (*Room, error)
	Delete(hash string) error
	Save(item *Room) error
	List(limit int) ([]*Room, error)
}
