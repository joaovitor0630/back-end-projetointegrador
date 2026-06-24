package database

import (
	"database/sql"
	"fmt"
	"plataforma-flamboyant/internal/models"
)

func CriarUsuario(u models.Usuario) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO usuarios (nome, departamento, cargo, perfil_id) VALUES ($1, $2, $3, $4) RETURNING usuario_id",
		u.Nome, u.Departamento, u.Cargo, u.PerfilID,
	).Scan(&id)
	return id, err
}

func ListarUsuarios() ([]models.Usuario, error) {
	rows, err := DB.Query("SELECT usuario_id, nome, departamento, cargo, perfil_id FROM usuarios ORDER BY usuario_id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var u models.Usuario
		if err := rows.Scan(&u.UsuarioID, &u.Nome, &u.Departamento, &u.Cargo, &u.PerfilID); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return usuarios, nil
}

func BuscarUsuarioPorID(id int) (models.Usuario, error) {
	var u models.Usuario
	err := DB.QueryRow(
		"SELECT usuario_id, nome, departamento, cargo, perfil_id FROM usuarios WHERE usuario_id = $1", id,
	).Scan(&u.UsuarioID, &u.Nome, &u.Departamento, &u.Cargo, &u.PerfilID)
	return u, err
}

func AtualizarUsuario(id int, u models.Usuario) error {
	res, err := DB.Exec(
		"UPDATE usuarios SET nome = $1, departamento = $2, cargo = $3, perfil_id = $4 WHERE usuario_id = $5",
		u.Nome, u.Departamento, u.Cargo, u.PerfilID, id,
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

func DeletarUsuario(id int) error {
	// Verificar se o usuário tem vistorias vinculadas
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM vistorias WHERE usuario_id = $1", id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("não é possível excluir o usuário: existem %d vistoria(s) vinculada(s)", count)
	}

	res, err := DB.Exec("DELETE FROM usuarios WHERE usuario_id = $1", id)
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
