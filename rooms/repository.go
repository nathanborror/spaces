package rooms

import "github.com/nathanborror/gommon/auth"

// RoomRepository holds all the methods needed to save, delete, load and list User objects.
type RoomRepository interface {
	Load(hash string) (*Room, error)
	Delete(hash string) error
	Save(room *Room) error
	List(limit int) ([]*Room, error)
	ListMembers(hash string) ([]*auth.User, error)
}

// RoomMemberRepository defines methods for finding rooms users are in
type RoomMemberRepository interface {
	Load(room string, user string) (*RoomMember, error)
	Save(roomMember *RoomMember) error
	List(user string, limit int) ([]*Room, error)
	ListJoinable(user string, limit int) ([]*Room, error)
	Delete(hash string) error
}
