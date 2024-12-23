package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {

		r.Get("/health", app.healthCheckHandler)

		r.Route("/book", func(r chi.Router) {

			r.Post("/upload", app.bookUploadHandler)
			r.Get("/{bookID}", app.bookGetHandler)
		})
	})

	return r
}

func (app *Application) run(mux http.Handler) error {

	srv := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has already been started at %s", app.Config.Addr)

	return srv.ListenAndServe()
}
