package dropbox

// SharedFolder represents information about a Dropbox shared folder
type SharedFolder struct {
	ID          string `json:"id,omitempty"`
	Path        string `json:"path,omitempty"`
	AccessLevel string `json:"access_level,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Members     struct {
		User struct {
			UID         int    `json:"uid,omitempty"`
			DisplayName string `json:"display_name,omitempty"`
			SameTeam    bool   `json:"same_team"`
		} `json:"user"`
		Role   string `json:"role,omitempty"`
		Active bool   `json:"active"`
	} `json:"members"`
}
