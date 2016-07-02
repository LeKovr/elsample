package app

import (
	"fmt"
	json "github.com/gorilla/rpc/v2/json2"
	"net/http"
	"time"

	auth "github.com/LeKovr/elsa-auth/psw"

	"github.com/LeKovr/go-base/database"
	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------

// ListArgs - аргументы метода List
type ListArgs struct {
	Offset, By    int
	Phone, IP     string
	Before, After time.Time
}

// Record - таблица журнала авторизаций
type Record struct {
	Stamp         time.Time `xorm:"pk created"`
	IP            string    `xorm:"ip pk"`
	Phone, Status string
}

// Records - массив строк журнала
type Records []Record

// -----------------------------------------------------------------------------

// App - Класс сервера API
type App struct {
	DB   *database.DB
	Log  *logger.Log
	Auth *auth.App
}

// -----------------------------------------------------------------------------

// New - Конструктор сервера API
func New(db *database.DB, log *logger.Log, auth *auth.App) *App {

	a := App{Log: log.WithField("in", "app"), DB: db, Auth: auth}
	a.initDB()
	return &a
}

// -----------------------------------------------------------------------------

// List - выборка строк из журнала
func (a *App) List(r *http.Request, args *ListArgs, reply *Records) error {
	_, err := a.Auth.ParseJWT(r)
	if err != nil {
		a.Log.Errorf("JWT parse error: %+v", err)
		return &json.Error{Code: -32011, Message: "Auth required"}
	}

	// debug.Log.Printf("Called Records: %+v", args)

	if args.Before.IsZero() {
		args.Before = time.Now()
	}

	var recs Records
	err = a.DB.Engine.
		Where("phone like ?", args.Phone+"%").
		And("ip like ?", args.IP+"%").
		And("stamp > ?", args.After).
		And("stamp <= ?", args.Before).
		Limit(args.By, args.Offset).
		Asc("stamp", "ip").
		Find(&recs)
	if err != nil {
		a.Log.Errorf("Fetch records error: %+v", err)
		return &json.Error{Code: -32012, Message: "Fetch error"}
	}

	*reply = recs
	return nil
}

// -----------------------------------------------------------------------------

// initDB prepares database
func (a *App) initDB() {

	engine := a.DB.Engine

	err := engine.Sync(new(Record))
	if err != nil {
		a.Log.Fatalf("DB sync error: %v", err)
	}

	isempty, err := engine.IsTableEmpty("record")
	if err != nil {
		a.Log.Fatalf("DB checkempty error: %v", err)
	}

	if isempty {
		// fill demo data
		for i := 1; i < 255; i++ {
			r := Record{Phone: "89181234567", IP: fmt.Sprintf("127.0.0.%d", i), Status: "SUCCESS", Stamp: time.Now().Add(time.Duration(-1*i) * time.Minute)}
			if _, err := engine.Insert(&r); err != nil {
				a.Log.Debugf("Record add error: %+v", err)
			}
		}
	}
	return
}
