package requisicao

import (
	"encoding/json"
	"log"
	"net/http"
)

func ProcessarRequisicao(w http.ResponseWriter, r *http.Request, data any) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		log.Printf("erro ao decodificar requisição: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Corpo da requisição inválido: " + err.Error()})
		return err
	}
	return nil
}
