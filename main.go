package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/comail/colog"
	"github.com/jessevdk/go-flags"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"gopkg.in/tylerb/graceful.v1"

	"github.com/LeKovr/go-base/database"
	"github.com/LeKovr/go-base/jwtutil"

	"github.com/LeKovr/elsa-auth/psw/api/admin"
	"github.com/LeKovr/elsa-auth/psw/api/anon"
	"github.com/LeKovr/elsa-auth/psw/api/user"
	"github.com/LeKovr/elsa-auth/psw/mw/acl"
	"github.com/LeKovr/elsa-auth/psw/mw/jwt"

	//	"github.com/LeKovr/elsa/api/hello"

	"github.com/LeKovr/elsa/mw/ace"
	"github.com/LeKovr/elsa/mw/flow"
	"github.com/LeKovr/elsa/mw/logger"
	"github.com/LeKovr/elsa/mw/realip"
	"github.com/LeKovr/elsa/mw/render"
	"github.com/LeKovr/elsa/mw/rpc"
	"github.com/LeKovr/elsa/mw/sample"
	"github.com/LeKovr/elsa/mw/stats"

	"github.com/LeKovr/elsample/api"
)

// -----------------------------------------------------------------------------

// Flags defines local application flags
type Flags struct {
	Addr     string   `long:"http_addr"   default:"localhost:8080" description:"Http listen address"`
	Static   string   `long:"static"      default:"static"         description:"Static files dir"`
	LogLevel string   `long:"log_level"   default:"warn"           description:"Log level [warn|info|debug]"`
	Origins  []string `long:"http_origin"   description:"Allowed http origin(s)"`
	Version  bool     `long:"version"       description:"Show version and exit"`
}

// Config defines all of application flags
type Config struct {
	Flags
	Ace    ace.Flags      `group:"Ace Options"`
	Render render.Flags   `group:"Render Options"`
	Sample sample.Flags   `group:"Sample Options"`
	Stats  stats.Flags    `group:"Stats Options"`
	JWT    jwtutil.Flags  `group:"JWT Options"`
	JWTMW  jwt.Flags      `group:"Auth token Options"`
	DB     database.Flags `group:"Database Options"`
	Anon   anon.Flags     `group:"Password login Options"`
}

// -----------------------------------------------------------------------------

const (
	ctxRealIP  = "real-ip" // Context field for Real ip
	ctxSession = "session" // Context field for Session data
)

// -----------------------------------------------------------------------------

func main() {

	var cfg Config
	lg, _ := setUp(&cfg)

	lg.Printf("info: %s v %s. ELSA Sample Web app", path.Base(os.Args[0]), Version)
	lg.Print("info: Copyright (C) 2016, Alexey Kovrizhkin <ak@elfire.ru>")

	n := negroni.New()

	CORS := cors.New(cors.Options{
		AllowedOrigins: cfg.Origins,
	})
	token, _ := jwtutil.New(lg, &cfg.JWT)

	// Setup database
	db, err := database.New(cfg.DB.Driver, cfg.DB.Connect, database.Debug(cfg.DB.Debug))
	panicIfError(err)
	anonAPI, err := anon.New(lg, &cfg.Anon, db, token, ctxRealIP)
	panicIfError(err)
	adminAPI, err := admin.New(lg, db, ctxSession, ctxRealIP)
	panicIfError(err)

	Use(n,
		realip.New(lg, ctxRealIP), // used in logger
		logger.New(lg, ctxRealIP),
		negroni.NewRecovery(),
		flow.New(lg),
		// hideNames(lg, ".")
		jwt.New(lg, &cfg.JWTMW, token, ctxSession), // Load JWT token
		acl.New(lg, ctxSession, "/admin", "admin"),
		acl.New(lg, ctxSession, "/user", "user"),
		acl.New(lg, ctxSession, "/my", "user", "admin"),
	)
	UseIfAllowed(n, lg,
		stats.New(lg, &cfg.Stats),

		CORS,
		rpc.New(lg, "/api/v1", anonAPI),
		rpc.New(lg, "/user/api/v1", user.New(lg, db, ctxSession, ctxRealIP)),
		rpc.New(lg, "/my/api/v1", api.New(lg, db)),
		rpc.New(lg, "/admin/api/v1", adminAPI),
		//		rpc.New(lg, "/admin/api/v1", hello.New(lg, "admin area")),

		negroni.NewStatic(http.Dir(cfg.Static)),
		negroni.NewStatic(http.Dir("vendor")),
	)
	Use(n,
		render.New(lg, &cfg.Render, false),
		ace.New(lg, &cfg.Ace, true),
	)

	lg.Printf("info: Listening in %s", cfg.Addr)
	srv := &graceful.Server{
		Timeout: 5 * time.Second,
		Server:  &http.Server{Addr: cfg.Addr, Handler: n},
		ShutdownInitiated: func() {
			lg.Printf("info: Server is shutting down")
		},
	}
	srv.ListenAndServe()
}

// -----------------------------------------------------------------------------

func setUp(cfg *Config) (lg *log.Logger, err error) {

	p := flags.NewParser(cfg, flags.Default)

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

	lvl, err := colog.ParseLevel(cfg.LogLevel)
	panicIfError(err)

	colog.Register()
	colog.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	colog.SetMinLevel(lvl) //Warning)
	colog.SetDefaultLevel(lvl)

	cl := colog.NewCoLog(os.Stderr, "", log.Lshortfile|log.Ldate|log.Ltime)
	cl.SetMinLevel(lvl)
	cl.SetDefaultLevel(lvl)
	lg = cl.NewLogger() // same as logger := log.New(cl, "", 0)

	return
}

// -----------------------------------------------------------------------------

// Use calls negroni.Use for a slice of handlers
func Use(n *negroni.Negroni, handlers ...negroni.Handler) {
	for _, h := range handlers {
		n.Use(h)
	}
}

// -----------------------------------------------------------------------------

// UseIfAllowed calls midleware only if flow's prohibited flag is not set
func UseIfAllowed(n *negroni.Negroni, lg *log.Logger, handlers ...negroni.Handler) {
	for _, h := range handlers {
		n.Use(flow.NewHandler(lg, h))
	}
}

// -----------------------------------------------------------------------------

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
