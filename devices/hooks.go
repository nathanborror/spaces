package devices

import (
	"time"

	"github.com/jmoiron/modl"
)

// PreInsert sets the Created and Modified time before Device is saved.
func (d *Device) PreInsert(modl.SqlExecutor) error {
	if d.Created.IsZero() {
		d.Created = time.Now()
	}
	if d.Modified.IsZero() {
		d.Modified = time.Now()
	}
	return nil
}

// PreUpdate updates the Modified time before Device is updated.
func (d *Device) PreUpdate(modl.SqlExecutor) error {
	d.Modified = time.Now()
	return nil
}
