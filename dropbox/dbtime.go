package dropbox

import (
	"encoding/json"
	"time"
)

// DBTime allow marshalling and unmarshalling of time.
type DBTime time.Time

// UnmarshalJSON unmarshals a time according to the Dropbox format.
func (dbt *DBTime) UnmarshalJSON(data []byte) error {
	var s string
	var err error
	var t time.Time

	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}
	if t, err = time.ParseInLocation(DateFormat, s, time.UTC); err != nil {
		return err
	}
	if t.IsZero() {
		*dbt = DBTime(time.Time{})
	} else {
		*dbt = DBTime(t)
	}
	return nil
}

// MarshalJSON marshals a time according to the Dropbox format.
func (dbt DBTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dbt).Format(DateFormat))
}
