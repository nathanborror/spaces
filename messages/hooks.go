package messages

import (
  "time"

  "github.com/jmoiron/modl"
)

// PreInsert sets the Created and Modified time before Message is saved.
func (m *Message) PreInsert(modl.SqlExecutor) error {
  if m.Created.IsZero() {
    m.Created = time.Now()
  }
  if m.Modified.IsZero() {
    m.Modified = time.Now()
  }
  return nil
}

// PreUpdate updates the Modified time before Message is updated.
func (m *Message) PreUpdate(modl.SqlExecutor) error {
  m.Modified = time.Now()
  return nil
}
