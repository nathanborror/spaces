package devices

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/anachronistic/apns"
	"github.com/jmoiron/modl"
	_ "github.com/mattn/go-sqlite3"
)

type sqlDeviceRepository struct {
	dbmap *modl.DbMap
}

// DeviceSQLRepository returns a new sqlDeviceRepository or panics if it cannot
func DeviceSQLRepository(filename string) DeviceRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic("Error connecting to db: " + err.Error())
	}
	r := &sqlDeviceRepository{
		dbmap: modl.NewDbMap(db, modl.SqliteDialect{}),
	}
	r.dbmap.TraceOn("", log.New(os.Stdout, "db: ", log.Lmicroseconds))
	r.dbmap.AddTable(Device{}).SetKeys(false, "token")
	r.dbmap.CreateTablesIfNotExists()
	return r
}

func (r *sqlDeviceRepository) Save(device *Device) error {
	n, err := r.dbmap.Update(device)
	if err != nil {
		panic(err)
	}
	if n == 0 {
		err = r.dbmap.Insert(device)
	}
	return err
}

func (r *sqlDeviceRepository) Delete(hash string) error {
	_, err := r.dbmap.Exec("DELETE FROM device WHERE hash=?", hash)
	return err
}

func (r *sqlDeviceRepository) Load(user string) ([]*Device, error) {
	obj := []*Device{}
	err := r.dbmap.Select(&obj, "SELECT * FROM device WHERE user = ?", user)
	return obj, err
}

func (r *sqlDeviceRepository) List(users []string) ([]*Device, error) {
	obj := []*Device{}
	err := r.dbmap.Select(&obj, "SELECT * FROM device WHERE user IN (?)", strings.Join(users, ", "))
	return obj, err
}

func (r *sqlDeviceRepository) Push(users []string, message string) error {
	devices := []*Device{}
	err := r.dbmap.Select(&devices, "SELECT * FROM device WHERE user IN (?)", strings.Join(users, ", "))
	if err != nil {
		return err
	}

	payload := apns.NewPayload()
	payload.Alert = message
	payload.Badge = 1 // TODO: Make this more accurate
	payload.Sound = "bingbong.aiff"

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", "SpacesCert.pem", "SpacesKeyNoEnc.pem")

	for _, device := range devices {
		pn := apns.NewPushNotification()
		pn.DeviceToken = device.Token
		pn.AddPayload(payload)
		resp := client.Send(pn)

		alert, _ := pn.PayloadString()
		if resp.Error != nil {
			log.Println("APNS Error: ", resp.Error)
		} else {
			log.Println("APNS Alert: ", alert)
			log.Println("APNS Success: ", resp.Success)
		}
	}

	return nil
}
