package main

import (
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"web.beckerlabs.dev/internal/models"
)

type application struct {
	logger        *slog.Logger
	posts         models.PostsModelInterface
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		posts:         &models.PostsModel{},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Starting server", "addr", *addr)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
