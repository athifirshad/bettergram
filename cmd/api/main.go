package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"athifirshad.com/bettergram/internal/data"
	"athifirshad.com/bettergram/internal/database"

	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL  string
	httpPort int
	db       struct {
		dsn string
	}
}

type application struct {
	config config
	db     *database.DB
	logger *slog.Logger
	wg     sync.WaitGroup
	data   data.Models
}

func run(logger *slog.Logger) error {
	var cfg config

	flag.StringVar(&cfg.baseURL, "base-url", "http://localhost:4000", "base URL for the application")
	flag.IntVar(&cfg.httpPort, "http-port", 4000, "port to listen on for HTTP requests")

	flag.Parse()

	cfg.db.dsn = os.Getenv("DB_DSN")
	if cfg.db.dsn == "" {
		return fmt.Errorf("DB_DSN environment variable is not set")
	}

	db, err := database.New(cfg.db.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	app := &application{
		config: cfg,
		db:     db,
		logger: logger,
		data:   data.NewModels(db.Pool),
	}

	return app.serveHTTP()
}
