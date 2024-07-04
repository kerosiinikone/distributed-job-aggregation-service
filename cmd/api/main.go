package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct{
	port 		int
	maxWorkers 	int
}

type application struct {
	logger 	*log.Logger
	cfg 	config
}

func main() {
	var cfg config
	
	flag.IntVar(&cfg.port, "port", 3000, "listen port")
	flag.IntVar(&cfg.maxWorkers, "maxWorkers", 4, "number of workers max")
	
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

	app := &application{
		logger: logger,
		cfg: cfg,
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
}