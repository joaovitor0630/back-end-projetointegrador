package database

import (
	"database/sql"
	"plataforma-flamboyant/internal/models"
)

func CriarLoja(loja models.Loja) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO lojas (nome, luc, segmento) VALUES ($1, $2, $3) RETURNING loja_id",
		loja.Nome, loja.LUC, loja.Segmento,
	).Scan(&id)
	return id, err
}

func ListarLojas() ([]models.Loja, error) {
	rows, err := DB.Query("SELECT loja_id, nome, luc, segmento, data_registro FROM lojas ORDER BY loja_id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lojas []models.Loja
	for rows.Next() {
		var l models.Loja
		if err := rows.Scan(&l.LojaID, &l.Nome, &l.LUC, &l.Segmento, &l.DataRegistro); err != nil {
			return nil, err
		}
		lojas = append(lojas, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lojas, nil
}

func BuscarLojaPorID(id int) (models.Loja, error) {
	var l models.Loja
	err := DB.QueryRow(
		"SELECT loja_id, nome, luc, segmento, data_registro FROM lojas WHERE loja_id = $1", id,
	).Scan(&l.LojaID, &l.Nome, &l.LUC, &l.Segmento, &l.DataRegistro)
	return l, err
}

func AtualizarLoja(id int, loja models.Loja) error {
	res, err := DB.Exec(
		"UPDATE lojas SET nome = $1, luc = $2, segmento = $3 WHERE loja_id = $4",
		loja.Nome, loja.LUC, loja.Segmento, id,
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

func DeletarLoja(id int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Deletar anexos de vistorias desta loja
	_, err = tx.Exec("DELETE FROM anexos WHERE vistoria_id IN (SELECT vistoria_id FROM vistorias WHERE loja_id = $1)", id)
	if err != nil {
		return err
	}
	// Deletar vistorias desta loja
	_, err = tx.Exec("DELETE FROM vistorias WHERE loja_id = $1", id)
	if err != nil {
		return err
	}
	// Deletar anexos de multas vinculadas a ocorrências desta loja
	_, err = tx.Exec("DELETE FROM anexos WHERE multa_id IN (SELECT multa_id FROM multas WHERE ocorrencia_id IN (SELECT ocorrencia_id FROM ocorrencias WHERE loja_id = $1))", id)
	if err != nil {
		return err
	}
	// Deletar multas vinculadas a ocorrências desta loja
	_, err = tx.Exec("DELETE FROM multas WHERE ocorrencia_id IN (SELECT ocorrencia_id FROM ocorrencias WHERE loja_id = $1)", id)
	if err != nil {
		return err
	}
	// Deletar anexos de ocorrências desta loja
	_, err = tx.Exec("DELETE FROM anexos WHERE ocorrencia_id IN (SELECT ocorrencia_id FROM ocorrencias WHERE loja_id = $1)", id)
	if err != nil {
		return err
	}
	// Deletar ocorrências desta loja
	_, err = tx.Exec("DELETE FROM ocorrencias WHERE loja_id = $1", id)
	if err != nil {
		return err
	}
	// Limpar auditoria órfã para vistorias, multas e ocorrências da loja
	_, err = tx.Exec("DELETE FROM auditoria WHERE entidade_tipo = 'vistoria' AND entidade_id IN (SELECT vistoria_id FROM vistorias WHERE loja_id = $1)", id)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM auditoria WHERE entidade_tipo = 'multa' AND entidade_id IN (SELECT multa_id FROM multas WHERE ocorrencia_id IN (SELECT ocorrencia_id FROM ocorrencias WHERE loja_id = $1))", id)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM auditoria WHERE entidade_tipo = 'ocorrencia' AND entidade_id IN (SELECT ocorrencia_id FROM ocorrencias WHERE loja_id = $1)", id)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM auditoria WHERE entidade_tipo = 'loja' AND entidade_id = $1", id)
	if err != nil {
		return err
	}

	// Deletar a loja
	res, err := tx.Exec("DELETE FROM lojas WHERE loja_id = $1", id)
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
