package boards

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/modl"
	_ "github.com/mattn/go-sqlite3"
)

type sqlBoardRepository struct {
	dbmap *modl.DbMap
}

// BoardSQLRepository returns a new sqlBoardRepository or panics if it cannot
func BoardSQLRepository(filename string) BoardRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	r := &sqlBoardRepository{
		dbmap: modl.NewDbMap(db, modl.SqliteDialect{}),
	}
	r.dbmap.TraceOn("", log.New(os.Stdout, "db: ", log.Lmicroseconds))
	r.dbmap.AddTable(Board{}).SetKeys(false, "hash")
	r.dbmap.CreateTablesIfNotExists()
	return r
}

func (r *sqlBoardRepository) Load(hash string) (*Board, error) {
	obj := []*Board{}
	err := r.dbmap.Select(&obj, "SELECT * FROM board WHERE hash=?", hash)
	if len(obj) != 1 {
		return nil, fmt.Errorf("expected 1 object, got %d", len(obj))
	}
	return obj[0], err
}

func (r *sqlBoardRepository) Save(board *Board) error {
	n, err := r.dbmap.Update(board)
	if err != nil {
		panic(err)
	}
	if n == 0 {
		err = r.dbmap.Insert(board)
	}
	return err
}

func (r *sqlBoardRepository) List(limit int) ([]*Board, error) {
	obj := []*Board{}
	err := r.dbmap.Select(&obj, "SELECT * FROM board ORDER BY created DESC LIMIT ?", limit)
	return obj, err
}

func (r *sqlBoardRepository) ListForRoom(room string) ([]*Board, error) {
	obj := []*Board{}
	err := r.dbmap.Select(&obj, "SELECT * FROM board WHERE room = ? ORDER BY created DESC", room)
	return obj, err
}

// Paths

type sqlPathRepository struct {
	dbmap *modl.DbMap
}

// PathSQLRepository returns a new sqlBoardRepository or panics if it cannot
func PathSQLRepository(filename string) PathRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	r := &sqlPathRepository{
		dbmap: modl.NewDbMap(db, modl.SqliteDialect{}),
	}
	r.dbmap.TraceOn("", log.New(os.Stdout, "db: ", log.Lmicroseconds))
	r.dbmap.AddTable(Path{}).SetKeys(false, "hash")
	r.dbmap.CreateTablesIfNotExists()
	return r
}

func (r *sqlPathRepository) Save(p *Path) error {
	n, err := r.dbmap.Update(p)
	if err != nil {
		panic(err)
	}
	if n == 0 {
		err = r.dbmap.Insert(p)
	}
	return err
}

func (r *sqlPathRepository) List(board string) ([]*Path, error) {
	obj := []*Path{}
	err := r.dbmap.Select(&obj, "SELECT * FROM path WHERE board = ?", board)
	return obj, err
}

func (r *sqlPathRepository) DeleteAll(board string) error {
	_, err := r.dbmap.Exec("DELETE FROM path WHERE board=?", board)
	return err
}

func (r *sqlPathRepository) Delete(hash string) error {
	_, err := r.dbmap.Exec("DELETE FROM path WHERE hash=?", hash)
	return err
}
