package models

import "time"

// ─── PERFIL ──────────────────────────────────────────────────────────────────

type Perfil struct {
	PerfilID   int    `json:"perfil_id"`
	TipoPerfil string `json:"tipo_perfil"`
	Descricao  string `json:"descricao"`
}

// ─── LOJA ────────────────────────────────────────────────────────────────────

type Loja struct {
	LojaID       int       `json:"loja_id"`
	Nome         string    `json:"nome"`
	LUC          string    `json:"luc"`
	Segmento     string    `json:"segmento"`
	DataRegistro time.Time `json:"data_registro"`
}

// ─── USUARIO ─────────────────────────────────────────────────────────────────

type Usuario struct {
	UsuarioID    int    `json:"usuario_id"`
	PerfilID     *int   `json:"perfil_id"`
	Nome         string `json:"nome"`
	Departamento string `json:"departamento"`
	Cargo        string `json:"cargo"`
}

// ─── OCORRÊNCIA ──────────────────────────────────────────────────────────────

type AreaResponsavel string

const (
	AreaArq AreaResponsavel = "Arq"
	AreaEng AreaResponsavel = "Eng"
	AreaBri AreaResponsavel = "Bri"
)

func (a AreaResponsavel) IsValid() bool {
	switch a {
	case AreaArq, AreaEng, AreaBri:
		return true
	}
	return false
}

type CategoriaOcorrencia string

const (
	CatManutencao            CategoriaOcorrencia = "Manutenção"
	CatConservacao           CategoriaOcorrencia = "Conservação"
	CatConservacaoResiduos   CategoriaOcorrencia = "Conservação de resíduos"
	CatArquiteturaPaisagismo CategoriaOcorrencia = "Arquitetura e Paisagismo"
	CatSeguranca             CategoriaOcorrencia = "Segurança"
	CatBrigada               CategoriaOcorrencia = "Brigada"
	CatEngenharia            CategoriaOcorrencia = "Engenharia"
)

func (c CategoriaOcorrencia) IsValid() bool {
	switch c {
	case CatManutencao, CatConservacao, CatConservacaoResiduos, CatArquiteturaPaisagismo, CatSeguranca, CatBrigada, CatEngenharia:
		return true
	}
	return false
}

type Ocorrencia struct {
	OcorrenciaID    int                 `json:"ocorrencia_id"`
	LojaID          int                 `json:"loja_id"`
	AreaResponsavel AreaResponsavel     `json:"area_responsavel"`
	Categoria       CategoriaOcorrencia `json:"categoria"`
	Assunto         string              `json:"assunto"`
	Descricao       string              `json:"descricao"`
	DataRegistro    time.Time           `json:"data_registro"`
}

// ─── MULTA ───────────────────────────────────────────────────────────────────

type StatusMulta string

const (
	StatusPendente  StatusMulta = "Pendente"
	StatusFaturada  StatusMulta = "Faturada"
	StatusCancelada StatusMulta = "Cancelada"
)

func (s StatusMulta) IsValid() bool {
	switch s {
	case StatusPendente, StatusFaturada, StatusCancelada:
		return true
	}
	return false
}

type CategoriaMulta string

const (
	CatMultaManutencao  CategoriaMulta = "Manutenção e Engenharia"
	CatMultaArquitetura CategoriaMulta = "Arquitetura e Paisagismo"
	CatMultaConservacao CategoriaMulta = "Conservação e Resíduos"
	CatMultaSeguranca   CategoriaMulta = "Segurança e Brigada"
)

func (c CategoriaMulta) IsValid() bool {
	switch c {
	case CatMultaManutencao, CatMultaArquitetura, CatMultaConservacao, CatMultaSeguranca:
		return true
	}
	return false
}

type Multa struct {
	MultaID      int            `json:"multa_id"`
	OcorrenciaID int            `json:"ocorrencia_id"`
	Categoria    CategoriaMulta `json:"categoria"`
	Assunto      string         `json:"assunto"`
	ValorMulta   float64        `json:"valor_multa"`
	Status       StatusMulta    `json:"status"`
}

// ─── VISTORIA ─────────────────────────────────────────────────────────────────

type CategoriaVistoria string

const (
	CatReformarLoja CategoriaVistoria = "Reformar loja"
	CatNovoContrato CategoriaVistoria = "Novo Contrato"
)

func (c CategoriaVistoria) IsValid() bool {
	switch c {
	case CatReformarLoja, CatNovoContrato:
		return true
	}
	return false
}

type Vistoria struct {
	VistoriaID      int               `json:"vistoria_id"`
	LojaID          int               `json:"loja_id"`
	UsuarioID       int               `json:"usuario_id"`
	AreaResponsavel AreaResponsavel   `json:"area_responsavel"`
	Categoria       CategoriaVistoria `json:"categoria"`
	Assunto         *string           `json:"assunto"`
	Descricao       string            `json:"descricao"`
	DataRegistro    time.Time         `json:"data_registro"`
}

// ─── ANEXO ───────────────────────────────────────────────────────────────────

type Anexo struct {
	AnexoID      int    `json:"anexo_id"`
	NomeArquivo  string `json:"nome_arquivo"`
	TipoMime     string `json:"tipo_mime"`
	TamanhoBytes int    `json:"tamanho_bytes"`
	Conteudo     []byte `json:"conteudo,omitempty"`
	OcorrenciaID *int   `json:"ocorrencia_id,omitempty"`
	MultaID      *int   `json:"multa_id,omitempty"`
	VistoriaID   *int   `json:"vistoria_id,omitempty"`
}

// ─── AUDITORIA ───────────────────────────────────────────────────────────────

type Auditoria struct {
	AuditoriaID   int       `json:"auditoria_id"`
	EntidadeTipo  string    `json:"entidade_tipo"`
	EntidadeID    int       `json:"entidade_id"`
	Acao          string    `json:"acao"`
	Usuario       string    `json:"usuario"`
	Campo         *string   `json:"campo,omitempty"`
	ValorAnterior *string   `json:"valor_anterior,omitempty"`
	ValorNovo     *string   `json:"valor_novo,omitempty"`
	Detalhes      *string   `json:"detalhes,omitempty"`
	DataHora      time.Time `json:"data_hora"`
}
