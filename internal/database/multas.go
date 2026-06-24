package database

import (
	"database/sql"
	"plataforma-flamboyant/internal/models"
)

func CriarMulta(m models.Multa) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO multas (ocorrencia_id, categoria, assunto, valor_multa, status) VALUES ($1, $2, $3, $4, $5) RETURNING multa_id",
		m.OcorrenciaID, m.Categoria, m.Assunto, m.ValorMulta, m.Status,
	).Scan(&id)
	return id, err
}

func ListarMultas() ([]models.Multa, error) {
	rows, err := DB.Query("SELECT multa_id, ocorrencia_id, categoria, assunto, valor_multa, status FROM multas ORDER BY multa_id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Multa
	for rows.Next() {
		var m models.Multa
		if err := rows.Scan(&m.MultaID, &m.OcorrenciaID, &m.Categoria, &m.Assunto, &m.ValorMulta, &m.Status); err != nil {
			return nil, err
		}
		lista = append(lista, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lista, nil
}

func BuscarMultaPorID(id int) (models.Multa, error) {
	var m models.Multa
	err := DB.QueryRow(
		"SELECT multa_id, ocorrencia_id, categoria, assunto, valor_multa, status FROM multas WHERE multa_id = $1", id,
	).Scan(&m.MultaID, &m.OcorrenciaID, &m.Categoria, &m.Assunto, &m.ValorMulta, &m.Status)
	return m, err
}

func AtualizarMulta(id int, m models.Multa) error {
	res, err := DB.Exec(
		"UPDATE multas SET ocorrencia_id = $1, categoria = $2, assunto = $3, valor_multa = $4, status = $5 WHERE multa_id = $6",
		m.OcorrenciaID, m.Categoria, m.Assunto, m.ValorMulta, m.Status, id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeletarMulta(id int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Deletar anexos vinculados à multa
	_, err = tx.Exec("DELETE FROM anexos WHERE multa_id = $1", id)
	if err != nil {
		return err
	}
	// Limpar auditoria órfã
	_, err = tx.Exec("DELETE FROM auditoria WHERE entidade_tipo = 'multa' AND entidade_id = $1", id)
	if err != nil {
		return err
	}

	// Deletar a multa
	res, err := tx.Exec("DELETE FROM multas WHERE multa_id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}
