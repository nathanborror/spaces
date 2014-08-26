package dropbox

// Entry represents the metadata of a file or folder.
type Entry struct {
	Bytes       int     `json:"bytes,omitempty"`        // Size of the file in bytes.
	ClientMtime DBTime  `json:"client_mtime,omitempty"` // Modification time set by the client when added.
	Contents    []Entry `json:"contents,omitempty"`     // List of children for a directory.
	Hash        string  `json:"hash,omitempty"`         // Hash of this entry.
	Icon        string  `json:"icon,omitempty"`         // Name of the icon displayed for this entry.
	IsDeleted   bool    `json:"is_deleted,omitempty"`   // true if this entry was deleted.
	IsDir       bool    `json:"is_dir,omitempty"`       // true if this entry is a directory.
	MimeType    string  `json:"mime_type,omitempty"`    // MimeType of this entry.
	Modified    DBTime  `json:"modified,omitempty"`     // Date of last modification.
	Path        string  `json:"path,omitempty"`         // Absolute path of this entry.
	Revision    string  `json:"rev,omitempty"`          // Unique ID for this file revision.
	Root        string  `json:"root,omitempty"`         // dropbox or sandbox.
	Size        string  `json:"size,omitempty"`         // Size of the file humanized/localized.
	ThumbExists bool    `json:"thumb_exists,omitempty"` // true if a thumbnail is available for this entry.
}
