package boards

import (
	"time"

	"github.com/jmoiron/modl"
)

// PreInsert sets the Created and Modified time before Board is saved.
func (b *Board) PreInsert(modl.SqlExecutor) error {
	if b.Created.IsZero() {
		b.Created = time.Now()
	}
	if b.Modified.IsZero() {
		b.Modified = time.Now()
	}
	return nil
}

// PreUpdate updates the Modified time before Board is updated.
func (b *Board) PreUpdate(modl.SqlExecutor) error {
	b.Modified = time.Now()
	return nil
}
