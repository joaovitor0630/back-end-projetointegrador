package database

import (
	"database/sql"
	"fmt"
	"plataforma-flamboyant/internal/models"
)

func CriarPerfil(perfil models.Perfil) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO perfil (tipo_perfil, descricao) VALUES ($1, $2) RETURNING perfil_id",
		perfil.TipoPerfil, perfil.Descricao,
	).Scan(&id)
	return id, err
}

func ListarPerfis() ([]models.Perfil, error) {
	rows, err := DB.Query("SELECT perfil_id, tipo_perfil, descricao FROM perfil ORDER BY perfil_id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perfis []models.Perfil
	for rows.Next() {
		var p models.Perfil
		if err := rows.Scan(&p.PerfilID, &p.TipoPerfil, &p.Descricao); err != nil {
			return nil, err
		}
		perfis = append(perfis, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return perfis, nil
}

func BuscarPerfilPorID(id int) (models.Perfil, error) {
	var p models.Perfil
	err := DB.QueryRow(
		"SELECT perfil_id, tipo_perfil, descricao FROM perfil WHERE perfil_id = $1", id,
	).Scan(&p.PerfilID, &p.TipoPerfil, &p.Descricao)
	return p, err
}

func AtualizarPerfil(id int, p models.Perfil) error {
	res, err := DB.Exec(
		"UPDATE perfil SET tipo_perfil = $1, descricao = $2 WHERE perfil_id = $3",
		p.TipoPerfil, p.Descricao, id,
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

func DeletarPerfil(id int) error {
	// Verificar se há usuários vinculados a este perfil
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM usuarios WHERE perfil_id = $1", id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("não é possível excluir o perfil: existem %d usuário(s) vinculado(s)", count)
	}

	res, err := DB.Exec("DELETE FROM perfil WHERE perfil_id = $1", id)
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
