package database

import (
	"database/sql"
	"fmt"
	"plataforma-flamboyant/internal/models"
)

func CriarOcorrencia(o models.Ocorrencia) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO ocorrencias (loja_id, area_responsavel, categoria, assunto, descricao) VALUES ($1, $2, $3, $4, $5) RETURNING ocorrencia_id",
		o.LojaID, o.AreaResponsavel, o.Categoria, o.Assunto, o.Descricao,
	).Scan(&id)
	return id, err
}

func ListarOcorrencias() ([]models.Ocorrencia, error) {
	rows, err := DB.Query("SELECT ocorrencia_id, loja_id, area_responsavel, categoria, assunto, descricao, data_registro FROM ocorrencias ORDER BY ocorrencia_id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Ocorrencia
	for rows.Next() {
		var o models.Ocorrencia
		if err := rows.Scan(&o.OcorrenciaID, &o.LojaID, &o.AreaResponsavel, &o.Categoria, &o.Assunto, &o.Descricao, &o.DataRegistro); err != nil {
			return nil, err
		}
		lista = append(lista, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lista, nil
}

func BuscarOcorrenciaPorID(id int) (models.Ocorrencia, error) {
	var o models.Ocorrencia
	err := DB.QueryRow(
		"SELECT ocorrencia_id, loja_id, area_responsavel, categoria, assunto, descricao, data_registro FROM ocorrencias WHERE ocorrencia_id = $1", id,
	).Scan(&o.OcorrenciaID, &o.LojaID, &o.AreaResponsavel, &o.Categoria, &o.Assunto, &o.Descricao, &o.DataRegistro)
	return o, err
}

func AtualizarOcorrencia(id int, o models.Ocorrencia) error {
	var oldLojaID int
	err := DB.QueryRow("SELECT loja_id FROM ocorrencias WHERE ocorrencia_id = $1", id).Scan(&oldLojaID)
	if err != nil {
		return err
	}
	if oldLojaID != o.LojaID {
		var count int
		err = DB.QueryRow("SELECT count(*) FROM multas WHERE ocorrencia_id = $1", id).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("não é possível alterar a loja pois há multas vinculadas a esta ocorrência")
		}
	}

	res, err := DB.Exec(
		"UPDATE ocorrencias SET loja_id = $1, area_responsavel = $2, categoria = $3, assunto = $4, descricao = $5 WHERE ocorrencia_id = $6",
		o.LojaID, o.AreaResponsavel, o.Categoria, o.Assunto, o.Descricao, id,
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

func DeletarOcorrencia(id int) error {
	var count int
	err := DB.QueryRow("SELECT count(*) FROM multas WHERE ocorrencia_id = $1", id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("não é possível excluir a ocorrência pois há multas ativas vinculadas a ela")
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Deletar anexos vinculados diretamente à ocorrência
	_, err = tx.Exec("DELETE FROM anexos WHERE ocorrencia_id = $1", id)
	if err != nil {
		return err
	}

	// Limpar auditoria
	_, err = tx.Exec("DELETE FROM auditoria WHERE entidade_tipo = 'ocorrencia' AND entidade_id = $1", id)
	if err != nil {
		return err
	}

	res, err := tx.Exec("DELETE FROM ocorrencias WHERE ocorrencia_id = $1", id)
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
