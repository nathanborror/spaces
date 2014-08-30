package devices

import (
	"net/http"

	"github.com/nathanborror/gommon/auth"
)

var repo = DeviceSQLRepository("db.sqlite3")
var userRepo = auth.AuthSQLRepository("db.sqlite3")

func check(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// SaveHandler saves a item
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetAuthenticatedUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	token := r.FormValue("token")
	make := r.FormValue("type")

	d := &Device{Token: token, Make: make, User: user.Hash}
	err = repo.Save(d)
	check(err, w)

	http.Redirect(w, r, "/", http.StatusFound)
}
