package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"gestao/internal/model"
	"gestao/internal/repository"
	"math"
	"time"
)

type ContaReceberService struct {
	repository *repository.Repository
	db         *sql.DB
}

func (s *ContaReceberService) CriarContaReceber(ctx context.Context, conta *model.ContaReceberCriar) error {
	if err := conta.Validar(); err != nil {
		return err
	}
	emissao, err := time.Parse("2006-01-02", conta.DtEmissao)
	if err != nil {
		return errors.New("data de emissão inválida")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Garante que cliente, categoria, condição e forma existam e que a forma faça parte da condição escolhida.
	var valido bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tb_clientes WHERE id=$1) AND EXISTS(SELECT 1 FROM tb_categorias_contas_receber WHERE id=$2) AND EXISTS(SELECT 1 FROM tb_condicao_forma_pagamento WHERE id_condicao=$3 AND id_forma_pagamento=$4)`, conta.IDCliente, conta.IDCategoria, conta.IDCondicaoPagamento, conta.IDFormaPagamento).Scan(&valido); err != nil {
		return err
	}
	if !valido {
		return errors.New("cliente, categoria, condição ou forma de pagamento inválidos")
	}
	condicao, err := s.repository.CondicoesPagamento.BuscarPorID(ctx, tx, conta.IDCondicaoPagamento)
	if err != nil {
		return err
	}
	valores := dividirEmParcelas(conta.ValorTotal, int(condicao.QtdParcelas))
	vencimentos := make([]string, len(valores))
	for i := range vencimentos {
		dias := condicao.DiasPrimeiroVenc + int64(i)*condicao.IntervaloParcelas
		vencimentos[i] = emissao.AddDate(0, 0, int(dias)).Format("2006-01-02")
	}
	grupo, err := novoUUID()
	if err != nil {
		return err
	}
	if err := s.repository.ContasReceber.CriarParcelas(ctx, tx, conta, grupo, valores, vencimentos); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *ContaReceberService) ListarContasReceber(ctx context.Context, busca, vencimento, status string) ([]*model.ContaReceber, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	resultado, err := s.repository.ContasReceber.Listar(ctx, tx, busca, vencimento, status)
	if err != nil {
		return nil, err
	}
	return resultado, tx.Commit()
}
func (s *ContaReceberService) ReceberConta(ctx context.Context, id int64, valor float64, idFormaPagamentoReal int64) error {
	if valor <= 0 {
		return errors.New("o valor do recebimento deve ser maior que zero")
	}
	if idFormaPagamentoReal <= 0 {
		return errors.New("a forma de pagamento utilizada é obrigatória")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	conta, err := s.repository.ContasReceber.BuscarPorID(ctx, tx, id)
	if err != nil {
		return err
	}
	if conta.Status == "PAGO" {
		return repository.CONTA_RECEBER_QUITADA
	}
	if valor > conta.SaldoRestante {
		return errors.New("o valor do recebimento não pode ser maior que o saldo restante")
	}
	var formaExiste bool
	if err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tb_formas_pagamento WHERE id = $1)`, idFormaPagamentoReal).Scan(&formaExiste); err != nil {
		return err
	}
	if !formaExiste {
		return errors.New("forma de pagamento utilizada inválida")
	}
	if err = s.repository.ContasReceber.Receber(ctx, tx, id, valor, idFormaPagamentoReal); err != nil {
		return err
	}
	return tx.Commit()
}

func dividirEmParcelas(valor float64, parcelas int) []float64 {
	totalCentavos := int64(math.Round(valor * 100))
	base := totalCentavos / int64(parcelas)
	resto := totalCentavos % int64(parcelas)
	valores := make([]float64, parcelas)
	for i := range valores {
		centavos := base
		if int64(i) < resto {
			centavos++
		}
		valores[i] = float64(centavos) / 100
	}
	return valores
}
func novoUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}
