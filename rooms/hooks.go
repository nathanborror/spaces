package rooms

import (
	"github.com/jmoiron/modl"
	"time"
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
