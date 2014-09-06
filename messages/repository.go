package messages

// MessageRepository holds all the methods needed to save, delete, load and list User objects.
type MessageRepository interface {
	Load(hash string) (*Message, error)
	Delete(hash string) error
	Save(message *Message) error
	List(room string, count int) (MessageList, error)
}
