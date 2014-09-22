package rooms

import "github.com/nathanborror/gommon/auth"

// RoomRepository holds all the methods needed to save, delete, load and list User objects.
type RoomRepository interface {
	Load(hash string) (*Room, error)
	LoadOneOnOne(user1 string, user2 string) (*Room, error)
	Delete(hash string) error
	Save(room *Room) error
	List(limit int) (RoomList, error)
}

// RoomMemberRepository defines methods for finding rooms users are in
type RoomMemberRepository interface {
	Load(room string, user string) (*RoomMember, error)
	Save(roomMember *RoomMember) error
	List(hash string) ([]*RoomMember, error)
	ListMembers(hash string) ([]*auth.User, error)
	ListRoomsForUser(user string, limit int) (RoomList, error)
	ListJoinableRoomsForUser(user string, limit int) (RoomList, error)
	Delete(hash string) error
}
