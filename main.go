package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/jessevdk/go-flags"

	"github.com/LeKovr/elsa"
	auth "github.com/LeKovr/elsa-auth/psw"
	"github.com/LeKovr/elsample/app"
	"github.com/LeKovr/go-base/database"
	"github.com/LeKovr/go-base/logger"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Version bool `long:"version" description:"Show version and exit"`
}

// Config defines all of application flags
type Config struct {
	Flags
	log    logger.Flags
	db     database.Flags
	auth   auth.Flags
	server elsa.Flags
}

// -----------------------------------------------------------------------------

func main() {

	var cfg Config
	db, log, _ := setUp(&cfg)
	defer log.Close()

	log.Infof("%s v %s. ELSA Sample Web app", path.Base(os.Args[0]), Version)
	log.Info("Copyright (C) 2016, Alexey Kovrizhkin <ak@elfire.ru>")

	s, _ := elsa.New(cfg.server.Addr, log, elsa.DB(db))
	s.Handle("/api", elsa.APIServer(s.RPC, log, cfg.server.Hosts...))

	appAuth, _ := auth.New(db, log, auth.Config(&cfg.auth))
	s.RPC.RegisterService(appAuth, "Auth")
	s.Handle("/auth", appAuth) // nginx auth subrequest

	appRecords := app.New(db, log, appAuth) // app.ParseJWT
	s.RPC.RegisterService(appRecords, "Records")

	s.RunServer()
}

// -----------------------------------------------------------------------------

func setUp(cfg *Config) (db *database.DB, log *logger.Log, err error) {

	p := flags.NewParser(nil, flags.Default)

	_, err = p.AddGroup("Application Options", "", cfg)
	panicIfError(err) // check Flags parse error

	_, err = p.AddGroup("Logging Options", "", &cfg.log)
	panicIfError(err)

	_, err = p.AddGroup("Database Options", "", &cfg.db)
	panicIfError(err)

	_, err = p.AddGroup("Auth package Options", "", &cfg.auth)
	panicIfError(err)

	_, err = p.AddGroup("Server Options", "", &cfg.server)
	panicIfError(err)

	_, err = p.Parse()
	if err != nil {
		os.Exit(1) // error message written already
	}

	if cfg.Version {
		// show version & exit
		fmt.Printf("%s\n%s\n%s", Version, Build, Commit)
		os.Exit(0)
	}

	// use all CPU cores for maximum performance
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create a new instance of the logger
	log, err = logger.New(logger.Dest(cfg.log.Dest), logger.Level(cfg.log.Level))
	panicIfError(err) // check Flags parse error

	// Setup database
	db, err = database.New(cfg.db.Driver, cfg.db.Connect, database.Debug(cfg.db.Debug))
	panicIfError(err) // check Flags parse error

	return
}

// -----------------------------------------------------------------------------

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
