package rooms

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/modl"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nathanborror/gommon/auth"
)

type sqlRoomRepository struct {
	dbmap *modl.DbMap
}

// RoomSQLRepository returns a new sqlRoomRepository or panics if it cannot
func RoomSQLRepository(filename string) RoomRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	r := &sqlRoomRepository{
		dbmap: modl.NewDbMap(db, modl.SqliteDialect{}),
	}
	r.dbmap.TraceOn("", log.New(os.Stdout, "db: ", log.Lmicroseconds))
	r.dbmap.AddTable(Room{}).SetKeys(false, "hash")
	r.dbmap.CreateTablesIfNotExists()
	return r
}

func (r *sqlRoomRepository) Load(hash string) (*Room, error) {
	obj := []*Room{}
	err := r.dbmap.Select(&obj, "SELECT * FROM room WHERE hash=?", hash)
	if len(obj) != 1 {
		return nil, fmt.Errorf("expected 1 object, got %d", len(obj))
	}
	return obj[0], err
}

func (r *sqlRoomRepository) Save(room *Room) error {
	n, err := r.dbmap.Update(room)
	if err != nil {
		panic(err)
	}
	if n == 0 {
		err = r.dbmap.Insert(room)
	}
	return err
}

func (r *sqlRoomRepository) Delete(hash string) error {
	_, err := r.dbmap.Exec("DELETE FROM room WHERE hash=?", hash)
	return err
}

func (r *sqlRoomRepository) List(limit int) ([]*Room, error) {
	obj := []*Room{}
	err := r.dbmap.Select(&obj, "SELECT * FROM room ORDER BY created DESC LIMIT ?", limit)
	return obj, err
}

func (r *sqlRoomRepository) ListMembers(hash string) ([]*auth.User, error) {
	obj := []*auth.User{}
	err := r.dbmap.Select(&obj, "SELECT * FROM user WHERE hash IN (SELECT user from roommember WHERE room = ?) ORDER BY created DESC", hash)
	return obj, err
}

// RoomMember

type sqlRoomMemberRepository struct {
	dbmap *modl.DbMap
}

// RoomMemberSQLRepository returns a new sqlRoomRepository or panics if it cannot
func RoomMemberSQLRepository(filename string) RoomMemberRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	r := &sqlRoomMemberRepository{
		dbmap: modl.NewDbMap(db, modl.SqliteDialect{}),
	}
	r.dbmap.TraceOn("", log.New(os.Stdout, "db: ", log.Lmicroseconds))
	r.dbmap.AddTable(RoomMember{}).SetKeys(false, "hash")
	r.dbmap.CreateTablesIfNotExists()
	return r
}

func (r *sqlRoomMemberRepository) Load(room string, user string) (*RoomMember, error) {
	obj := []*RoomMember{}
	err := r.dbmap.Select(&obj, "SELECT * FROM roommember WHERE room=? AND user=?", room, user)
	if len(obj) != 1 {
		return nil, fmt.Errorf("expected 1 object, got %d", len(obj))
	}
	return obj[0], err
}

func (r *sqlRoomMemberRepository) Save(rm *RoomMember) error {
	n, err := r.dbmap.Update(rm)
	if err != nil {
		panic(err)
	}
	if n == 0 {
		err = r.dbmap.Insert(rm)
	}
	return err
}

func (r *sqlRoomMemberRepository) List(user string, limit int) ([]*Room, error) {
	rooms := []*Room{}
	err := r.dbmap.Select(&rooms, "SELECT * FROM room WHERE hash IN (SELECT room FROM roommember WHERE user = ?) ORDER BY created DESC LIMIT ?", user, limit)
	return rooms, err
}

func (r *sqlRoomMemberRepository) Delete(hash string) error {
	_, err := r.dbmap.Exec("DELETE FROM roommember WHERE hash=?", hash)
	return err
}
