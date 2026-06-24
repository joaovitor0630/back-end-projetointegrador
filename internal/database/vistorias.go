package database

import (
	"database/sql"
	"plataforma-flamboyant/internal/models"
)

func CriarVistoria(v models.Vistoria) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO vistorias (loja_id, usuario_id, area_responsavel, categoria, assunto, descricao) VALUES ($1, $2, $3, $4, $5, $6) RETURNING vistoria_id",
		v.LojaID, v.UsuarioID, v.AreaResponsavel, v.Categoria, v.Assunto, v.Descricao,
	).Scan(&id)
	return id, err
}

func ListarVistorias() ([]models.Vistoria, error) {
	rows, err := DB.Query("SELECT vistoria_id, loja_id, usuario_id, area_responsavel, categoria, assunto, descricao, data_registro FROM vistorias ORDER BY vistoria_id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Vistoria
	for rows.Next() {
		var v models.Vistoria
		if err := rows.Scan(&v.VistoriaID, &v.LojaID, &v.UsuarioID, &v.AreaResponsavel, &v.Categoria, &v.Assunto, &v.Descricao, &v.DataRegistro); err != nil {
			return nil, err
		}
		lista = append(lista, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lista, nil
}

func BuscarVistoriaPorID(id int) (models.Vistoria, error) {
	var v models.Vistoria
	err := DB.QueryRow(
		"SELECT vistoria_id, loja_id, usuario_id, area_responsavel, categoria, assunto, descricao, data_registro FROM vistorias WHERE vistoria_id = $1", id,
	).Scan(&v.VistoriaID, &v.LojaID, &v.UsuarioID, &v.AreaResponsavel, &v.Categoria, &v.Assunto, &v.Descricao, &v.DataRegistro)
	return v, err
}

func AtualizarVistoria(id int, v models.Vistoria) error {
	res, err := DB.Exec(
		"UPDATE vistorias SET loja_id = $1, usuario_id = $2, area_responsavel = $3, categoria = $4, assunto = $5, descricao = $6 WHERE vistoria_id = $7",
		v.LojaID, v.UsuarioID, v.AreaResponsavel, v.Categoria, v.Assunto, v.Descricao, id,
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

func DeletarVistoria(id int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Deletar anexos vinculados à vistoria
	_, err = tx.Exec("DELETE FROM anexos WHERE vistoria_id = $1", id)
	if err != nil {
		return err
	}
	// Limpar auditoria órfã
	_, err = tx.Exec("DELETE FROM auditoria WHERE entidade_tipo = 'vistoria' AND entidade_id = $1", id)
	if err != nil {
		return err
	}

	// Deletar a vistoria
	res, err := tx.Exec("DELETE FROM vistorias WHERE vistoria_id = $1", id)
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
