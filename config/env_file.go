package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CarregarArquivoEnv carrega variáveis de um arquivo no formato .env/.envrc.
// Variáveis já fornecidas pelo ambiente têm precedência sobre o arquivo local.
func CarregarArquivoEnv(caminho string) error {
	arquivo, err := os.Open(caminho)
	if err != nil {
		return err
	}
	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)
	for linha := 1; scanner.Scan(); linha++ {
		texto := strings.TrimSpace(scanner.Text())
		if texto == "" || strings.HasPrefix(texto, "#") {
			continue
		}
		texto = strings.TrimSpace(strings.TrimPrefix(texto, "export "))

		chave, valor, encontrado := strings.Cut(texto, "=")
		chave = strings.TrimSpace(chave)
		if !encontrado || chave == "" {
			return fmt.Errorf("arquivo %s: linha %d inválida", caminho, linha)
		}
		if _, definido := os.LookupEnv(chave); definido {
			continue
		}

		valor = strings.Trim(strings.TrimSpace(valor), "\"'")
		if err := os.Setenv(chave, valor); err != nil {
			return fmt.Errorf("arquivo %s: não foi possível definir %s: %w", caminho, chave, err)
		}
	}

	return scanner.Err()
}
