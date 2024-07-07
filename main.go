package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anthdm/hollywood/actor"
	"github.com/joho/godotenv"
)

// EventStream ??

// Add different levels of actors -> Overseer, Scraper, Sender, etc

// Add a Pub/sub layer (later), where some actors constantly check
// for changes without externals prompts that trigger certain functionalities
// when finding new matches

type Config struct {
	Port       int
	MaxWorkers int
}

type Application struct {
	Logger *log.Logger
	Cfg    Config
	Engine *actor.Engine
	MPid   *actor.PID
}

func NewApplication(cfg Config, logger *log.Logger, engine *actor.Engine, pid *actor.PID) *Application {
	return &Application{
		Cfg:    cfg,
		Logger: logger,
		Engine: engine,
		MPid:   pid,
	}
}

func main() {
	var cfg Config

	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
  	}
	
	flag.IntVar(&cfg.Port, "port", 3000, "listen port")
	flag.IntVar(&cfg.MaxWorkers, "maxWorkers", 4, "number of workers max")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}	

	// https://github.com/anthdm/hollywood?tab=readme-ov-file#with-custom-configuration
	mPid := e.Spawn(NewManager(), "manager")


	app := NewApplication(cfg, logger, e, mPid)

	if err != nil {
		app.Logger.Fatalf("Error occurred in the Actor Engine: %v\n", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/find", app.findJobHandler)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	fmt.Printf("Listening on %s", srv.Addr)

	// Blocking
	srv.ListenAndServe()

	// Signal Completion -> Poison the manager / active Actors
}