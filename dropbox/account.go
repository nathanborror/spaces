package dropbox

// Token represents information about the oauth token
type Token struct {
	AccessToken string `json:"access_token"`
}

// Account represents information about the user account.
type Account struct {
	ReferralLink string `json:"referral_link,omitempty"` // URL for referral.
	DisplayName  string `json:"display_name,omitempty"`  // User name.
	UID          int    `json:"uid,omitempty"`           // User account ID.
	Country      string `json:"country,omitempty"`       // Country ISO code.
	QuotaInfo    struct {
		Shared int64 `json:"shared,omitempty"` // Quota for shared files.
		Quota  int64 `json:"quota,omitempty"`  // Quota in bytes.
		Normal int64 `json:"normal,omitempty"` // Quota for non-shared files.
	} `json:"quota_info"`
}
