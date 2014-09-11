package boards

// BoardRepository defines a board with bezier paths
type BoardRepository interface {
	Load(hash string) (*Board, error)
	Save(room *Board) error
	List(room string) ([]*Board, error)
	Clear(hash string) error
}

// PathRepository defines a bezier path for a board
type PathRepository interface {
	Save(path *Path) error
	List(board string) ([]*Path, error)
	Delete(hash string) error
}
