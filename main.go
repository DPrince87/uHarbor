package main

import (
	"flag"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
	DB     *gorm.DB
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "Server Port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	mux := app.routes()
	app.db()

	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		fmt.Println(err)
	}
}
