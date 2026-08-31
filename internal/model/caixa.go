package model

import "errors"

type Caixa struct {
	ID   int64  `json:"id"`
	Nome string `json:"nome"`
}
type ControleCaixa struct {
	ID            int64   `json:"id"`
	IDCaixa       int64   `json:"id_caixa"`
	Status        string  `json:"status"`
	ValorAbertura float64 `json:"valor_abertura"`
}
type FechamentoCaixa struct {
	ValorDinheiro float64 `json:"valor_dinheiro"`
	ValorCartao   float64 `json:"valor_cartao"`
	SenhaAcesso   string  `json:"senha_acesso"`
}

func (f FechamentoCaixa) Validar() error {
	if f.ValorDinheiro < 0 || f.ValorCartao < 0 {
		return errors.New("os valores de fechamento não podem ser negativos")
	}
	if f.SenhaAcesso == "" {
		return errors.New("a senha de acesso é obrigatória")
	}
	return nil
}
