package request

import "finalizacao-pedido-svc/internal/finalizar/models"

type FinalizarRequestBase struct {
	EmpresaID        *int8            `json:"empresa_id,omitempty"`
	CodigoCliente    *string          `json:"codigo_cliente,omitempty"`
	Email            *string          `json:"email,omitempty"`
	CPF              *string          `json:"cpf,omitempty"`
	Parcela          *int8            `json:"parcela,omitempty"`
	ValorTotalPedido *float64         `json:"valor_total_pedido,omitempty"`
	Observacao       *string          `json:"observacao,omitempty"`
	Frete            *float64         `json:"frete,omitempty"`
	Percentual       *float64         `json:"percentual,omitempty"`
	FormaPagamento   *string          `json:"forma_pagamento,omitempty"`
	Endereco         *string          `json:"endereco,omitempty"`
	Carrinho         *models.Carrinho `json:"carrinho,omitempty"`
}

type FinalizarRequestPix struct {
	FinalizarRequestBase
	TxID *string `json:"txid,omitempty"`
}

type FinalizarRequestCartaoCredito struct {
	FinalizarRequestBase
	CodigoAut *string `json:"codigo_aut,omitempty"`
	NSU       *string `json:"nsu,omitempty"`
	Bandeira  *string `json:"bandeira,omitempty"`
}

type FinalizarRequestBoleto struct {
	FinalizarRequestBase
	CodigoAut  *string `json:"codigo_aut,omitempty"`
	CodigoCond *string `json:"codigo_cond,omitempty"`
}
