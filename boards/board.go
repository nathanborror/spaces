package boards

import "time"

// Board defines a place to draw
type Board struct {
	Hash     string    `json:"hash"`
	Room     string    `json:"room"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// Path is a drawing path pertaining to a board
type Path struct {
	Hash  string `json:"hash"`
	Board string `json:"board"`
	Data  string `json:"data"`
	User  string `json:"user"`
}
