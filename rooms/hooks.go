package rooms

import (
	"time"

	"github.com/jmoiron/modl"
)

// PreInsert sets the Created and Modified time before Room is saved.
func (r *Room) PreInsert(modl.SqlExecutor) error {
	if r.Created.IsZero() {
		r.Created = time.Now()
	}
	if r.Modified.IsZero() {
		r.Modified = time.Now()
	}
	return nil
}

// PreUpdate updates the Modified time before Room is updated.
func (r *Room) PreUpdate(modl.SqlExecutor) error {
	r.Modified = time.Now()
	return nil
}

// RoomMember

// PreInsert sets the Created time before RoomMember is saved.
func (rm *RoomMember) PreInsert(modl.SqlExecutor) error {
	if rm.Created.IsZero() {
		rm.Created = time.Now()
	}
	return nil
}
