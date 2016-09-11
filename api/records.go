package api

import (
	"fmt"
	rpc "github.com/gorilla/rpc/v2/json2"
	"log"
	"net/http"
	"time"

	"github.com/LeKovr/go-base/database"
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

// Service holds service attributes
type Service struct {
	Log     *log.Logger
	DB      *database.DB
	Field   string // Context field for Sessiondata
	IPField string // `long:"logger_realip_field" default:"real-ip" description:"Context field for Real ip"`
}

// -----------------------------------------------------------------------------

// New - Конструктор сервера API
func New(logger *log.Logger, db *database.DB) *Service {
	srv := Service{Log: logger, DB: db}
	srv.initDB()
	return &srv
}

// -----------------------------------------------------------------------------

// List - выборка строк из журнала
func (srv *Service) List(r *http.Request, args *ListArgs, result *Records) error {

	// debug.Log.Printf("Called Records: %+v", args)

	if args.Before.IsZero() {
		args.Before = time.Now()
	}

	srv.Log.Printf("debug: called list(%+v)", args)

	var recs Records
	err := srv.DB.Engine.
		Where("phone like ?", "%"+args.Phone+"%").
		And("ip like ?", "%"+args.IP+"%").
		And("stamp > ?", args.After).
		And("stamp <= ?", args.Before).
		Limit(args.By, args.Offset).
		Asc("stamp", "ip").
		Find(&recs)
	if err != nil {
		srv.Log.Printf("error: Fetch records error: %+v", err)
		return &rpc.Error{Code: -32012, Message: "Fetch error"}
	}

	*result = recs
	return nil
}

// -----------------------------------------------------------------------------

// initDB prepares database
func (srv *Service) initDB() {

	engine := srv.DB.Engine

	err := engine.Sync(new(Record))
	if err != nil {
		srv.Log.Printf("fatal: DB sync error: %v", err)
	}

	isempty, err := engine.IsTableEmpty("record")
	if err != nil {
		srv.Log.Printf("fatal: DB checkempty error: %v", err)
	}

	if isempty {
		// fill demo data
		for i := 1; i < 255; i++ {
			r := Record{Phone: "89181234567", IP: fmt.Sprintf("127.0.0.%d", i), Status: "SUCCESS", Stamp: time.Now().Add(time.Duration(-1*i) * time.Minute)}
			if _, err := engine.Insert(&r); err != nil {
				srv.Log.Printf("debug: Record add error: %+v", err)
			}
		}
	}
	return
}
