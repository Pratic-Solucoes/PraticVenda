package model

import (
	"errors"
)

// Debito representa a entidade principal de débito no banco de dados
type Debito struct {
	ID              int64   `json:"id" db:"id"`
	IDFornecedor    int64   `json:"id_fornecedor" db:"id_fornecedor"`
	IDCategoria     *int64  `json:"id_categoria,omitempty" db:"id_categoria"`
	Descricao       string  `json:"descricao" db:"descricao"`
	NrDocumento     *string `json:"nr_documento,omitempty" db:"nr_documento"`
	NrNotaFiscal    *string `json:"nr_nota_fiscal,omitempty" db:"nr_nota_fiscal"`
	Valor           float64 `json:"valor" db:"valor"`
	DtEntrada       string  `json:"dt_entrada" db:"dt_entrada"`
	DtVencimento    string  `json:"dt_vencimento" db:"dt_vencimento"`
	NrParcela       int     `json:"nr_parcela" db:"nr_parcela"`
	NrTotalParcelas int     `json:"nr_total_parcelas" db:"nr_total_parcelas"`
	Status          string  `json:"status" db:"status"` // Ex: PENDENTE, PAGO, CANCELADO
	CreatedAt       string  `json:"created_at" db:"created_at"`
	UpdatedAt       string  `json:"updated_at" db:"updated_at"`

	// Relacionamentos opcionais
	Fornecedor *Fornecedor      `json:"fornecedor,omitempty"`
	Categoria  *CategoriaDebito `json:"categoria,omitempty"`
}

// CategoriaDebito representa a categoria do débito (baseado no SQL cria_tb_categorias_debito)
type CategoriaDebito struct {
	ID        int64  `json:"id" db:"id"`
	Nome      string `json:"nome" db:"nome"`
}

// DebitoAvulsoCriar é um DTO (Data Transfer Object) usado apenas para receber dados de criação
type DebitoAvulsoCriar struct {
	IDFornecedor    int64   `json:"id_fornecedor"`
	IDCategoria     *int64  `json:"id_categoria,omitempty"`
	Descricao       string  `json:"descricao"`
	NrDocumento     string  `json:"nr_documento,omitempty"`
	NrNotaFiscal    string  `json:"nr_nota_fiscal,omitempty"`
	Valor           float64 `json:"valor"`
	DtEntrada       string  `json:"dt_entrada"`    // Pode receber string "YYYY-MM-DD" no JSON
	DtVencimento    string  `json:"dt_vencimento"` // Pode receber string "YYYY-MM-DD" no JSON
	NrParcela       int     `json:"nr_parcela"`
	NrTotalParcelas int     `json:"nr_total_parcelas"`
}

// Validar verifica se os dados obrigatórios para criação foram preenchidos corretamente
func (d *DebitoAvulsoCriar) Validar() error {
	var erros []error
	if d.IDFornecedor == 0 {
		erros = append(erros, errors.New("o id do fornecedor não foi informado"))
	}
	if d.Descricao == "" {
		erros = append(erros, errors.New("a descrição não foi informada"))
	}
	if d.Valor <= 0 {
		erros = append(erros, errors.New("o valor deve ser maior que zero"))
	}
	if d.DtEntrada == "" {
		erros = append(erros, errors.New("a data de entrada não foi informada"))
	}
	if d.DtVencimento == "" {
		erros = append(erros, errors.New("a data de vencimento não foi informada"))
	}
	if d.NrParcela == 0 {
		erros = append(erros, errors.New("o número da parcela não foi informado"))
	}
	if d.NrTotalParcelas == 0 {
		erros = append(erros, errors.New("o número total de parcelas não foi informado"))
	}
	if len(erros) > 0 {
		return errors.Join(erros...)
	}
	return nil
}
