package database

import (
	"plataforma-flamboyant/internal/models"
)

func CriarAuditoria(a models.Auditoria) (int, error) {
	var id int
	err := DB.QueryRow(
		`INSERT INTO auditoria 
		(entidade_tipo, entidade_id, acao, usuario, campo, valor_anterior, valor_novo, detalhes) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING auditoria_id`,
		a.EntidadeTipo, a.EntidadeID, a.Acao, a.Usuario, a.Campo, a.ValorAnterior, a.ValorNovo, a.Detalhes,
	).Scan(&id)
	return id, err
}

func ListarAuditoriaPorEntidade(entidadeTipo string, entidadeID int) ([]models.Auditoria, error) {
	rows, err := DB.Query(
		`SELECT auditoria_id, entidade_tipo, entidade_id, acao, usuario, campo, valor_anterior, valor_novo, detalhes, data_hora 
		FROM auditoria 
		WHERE entidade_tipo = $1 AND entidade_id = $2 
		ORDER BY auditoria_id ASC`,
		entidadeTipo, entidadeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Auditoria
	for rows.Next() {
		var a models.Auditoria
		if err := rows.Scan(
			&a.AuditoriaID, &a.EntidadeTipo, &a.EntidadeID, &a.Acao, &a.Usuario,
			&a.Campo, &a.ValorAnterior, &a.ValorNovo, &a.Detalhes, &a.DataHora,
		); err != nil {
			return nil, err
		}
		lista = append(lista, a)
	}
	return lista, rows.Err()
}
