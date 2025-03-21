package models

import "time"

type Carrinho struct {
	CodCarrinho               string     `json:"cod_carrinho" gorm:"primaryKey;unique"`
	CodConta                  *int       `json:"cod_conta,omitempty"`
	CodCupom                  *int       `json:"cod_cupom,omitempty"`
	Status                    *string    `json:"status,omitempty"`
	CriadoEm                  *time.Time `json:"data_criacao,omitempty" gorm:"default:CURRENT_TIMESTAMP"`
	AtualizadoEm              *time.Time `json:"data_update,omitempty" gorm:"default:CURRENT_TIMESTAMP"`
	StatusFinalizacao         *int       `json:"status_finalizacao,omitempty" gorm:"default:0"`
	OrcamentoProtheus         *string    `json:"orcamento_protheus,omitempty"`
	OrcamentoFinalizado       int        `json:"orcamento_finalizado" gorm:"default:0"`
	ProcessadoProtheus        int        `json:"processado_protheus" gorm:"default:0"`
	ValorAtual                *float64   `json:"valor_atual,omitempty"`
	FreteAtual                *float64   `json:"frete_atual,omitempty"`
	ProcessadoCresceVendas    int        `json:"processado_cresce_vendas" gorm:"default:0"`
	AposProcessamentoProtheus int        `json:"apos_processamento_protheus" gorm:"default:0"`
	EmProcessamentoEcommerce  int        `json:"em_processamento_ecommerce" gorm:"default:0"`
}
