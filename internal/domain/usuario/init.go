package usuario

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func InicializarModulo(db *sql.DB) *chi.Mux {

	repository := NewRepository(db)
	service := NewService(repository, db)
	controller := NewController(service)

	return CarregarRotas(controller)
}
