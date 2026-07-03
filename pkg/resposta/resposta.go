package resposta

import (
	"encoding/json"
	"log"
	"net/http"
)

func Padrao(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("erro ao codificar resposta: %v", err)
		}
	}

}
