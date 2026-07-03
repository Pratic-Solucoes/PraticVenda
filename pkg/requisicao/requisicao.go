package requisicao

import (
	"encoding/json"
	"log"
	"net/http"
)

func ProcessarRequisicao(w http.ResponseWriter, r *http.Request, data any) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		log.Printf("erro ao decodificar requisição: %v", err)
		http.Error(w, "Erro ao processar requisição", http.StatusBadRequest)
		return err
	}
	return nil
}
