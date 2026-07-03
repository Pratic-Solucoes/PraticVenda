package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type application struct {
	API_port           string
	API_maxHeaderBytes int
	API_readTimeout    time.Duration
	API_writeTimeout   time.Duration
}

func (app *application) iniciarApp(r *chi.Mux) error {
	fmt.Printf("Servidor iniciado na porta %s\n", app.API_port)
	srv := &http.Server{
		Addr:           app.API_port,
		Handler:        r,
		ReadTimeout:    app.API_readTimeout,
		WriteTimeout:   app.API_writeTimeout,
		MaxHeaderBytes: app.API_maxHeaderBytes,
	}
	return srv.ListenAndServe()
}
