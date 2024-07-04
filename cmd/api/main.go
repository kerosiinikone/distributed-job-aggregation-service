package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anthdm/hollywood/actor"
	actors "github.com/kerosiinikone/go-actors-project/internal"
)

type config struct{
	port 		int
	maxWorkers 	int
}

type Application struct {
	Logger 	*log.Logger
	Cfg 	config
	Engine 	*actor.Engine
	MPid	*actor.PID
}

func NewApplication(cfg config, logger *log.Logger, engine *actor.Engine, pid *actor.PID) *Application {
	return &Application{
		Cfg: cfg,
		Logger: logger,
		Engine: engine,
		MPid: pid,
	}
}

func main() {
	var cfg config
	
	flag.IntVar(&cfg.port, "port", 3000, "listen port")
	flag.IntVar(&cfg.maxWorkers, "maxWorkers", 4, "number of workers max")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
		}	
	mPid := e.Spawn(actors.NewManager(), "manager") // opts ??

	app := NewApplication(cfg, logger, e, mPid)

	if err != nil {
		app.Logger.Fatalf("Error occurred in the Actor Engine: %v\n", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/find", app.handleFindJob)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: mux,
	}

	fmt.Printf("Listening on %s", srv.Addr)

	// Blocking
	srv.ListenAndServe()

	// Signal Completion -> Poison the manager / active Actors
}