package database

import (
	"database/sql"
	"plataforma-flamboyant/internal/models"
)

func CriarAnexo(a models.Anexo) (int, error) {
	var id int
	err := DB.QueryRow(
		"INSERT INTO anexos (nome_arquivo, tipo_mime, tamanho_bytes, conteudo, ocorrencia_id, multa_id, vistoria_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING anexo_id",
		a.NomeArquivo, a.TipoMime, a.TamanhoBytes, a.Conteudo, a.OcorrenciaID, a.MultaID, a.VistoriaID,
	).Scan(&id)
	return id, err
}

func ListarAnexosPorOcorrencia(id int) ([]models.Anexo, error) {
	rows, err := DB.Query("SELECT anexo_id, nome_arquivo, tipo_mime, tamanho_bytes, ocorrencia_id, multa_id, vistoria_id FROM anexos WHERE ocorrencia_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Anexo
	for rows.Next() {
		var a models.Anexo
		if err := rows.Scan(&a.AnexoID, &a.NomeArquivo, &a.TipoMime, &a.TamanhoBytes, &a.OcorrenciaID, &a.MultaID, &a.VistoriaID); err != nil {
			return nil, err
		}
		lista = append(lista, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lista, nil
}

func ListarAnexosPorMulta(id int) ([]models.Anexo, error) {
	rows, err := DB.Query("SELECT anexo_id, nome_arquivo, tipo_mime, tamanho_bytes, ocorrencia_id, multa_id, vistoria_id FROM anexos WHERE multa_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Anexo
	for rows.Next() {
		var a models.Anexo
		if err := rows.Scan(&a.AnexoID, &a.NomeArquivo, &a.TipoMime, &a.TamanhoBytes, &a.OcorrenciaID, &a.MultaID, &a.VistoriaID); err != nil {
			return nil, err
		}
		lista = append(lista, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lista, nil
}

func ListarAnexosPorVistoria(id int) ([]models.Anexo, error) {
	rows, err := DB.Query("SELECT anexo_id, nome_arquivo, tipo_mime, tamanho_bytes, ocorrencia_id, multa_id, vistoria_id FROM anexos WHERE vistoria_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Anexo
	for rows.Next() {
		var a models.Anexo
		if err := rows.Scan(&a.AnexoID, &a.NomeArquivo, &a.TipoMime, &a.TamanhoBytes, &a.OcorrenciaID, &a.MultaID, &a.VistoriaID); err != nil {
			return nil, err
		}
		lista = append(lista, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lista, nil
}

func BuscarAnexoPorID(id int) (models.Anexo, error) {
	var a models.Anexo
	err := DB.QueryRow(
		"SELECT anexo_id, nome_arquivo, tipo_mime, tamanho_bytes, conteudo, ocorrencia_id, multa_id, vistoria_id FROM anexos WHERE anexo_id = $1", id,
	).Scan(&a.AnexoID, &a.NomeArquivo, &a.TipoMime, &a.TamanhoBytes, &a.Conteudo, &a.OcorrenciaID, &a.MultaID, &a.VistoriaID)
	return a, err
}

func AtualizarAnexo(id int, a models.Anexo) error {
	res, err := DB.Exec(
		"UPDATE anexos SET nome_arquivo = $1, tipo_mime = $2, tamanho_bytes = $3, conteudo = $4, ocorrencia_id = $5, multa_id = $6, vistoria_id = $7 WHERE anexo_id = $8",
		a.NomeArquivo, a.TipoMime, a.TamanhoBytes, a.Conteudo, a.OcorrenciaID, a.MultaID, a.VistoriaID, id,
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

func DeletarAnexo(id int) error {
	res, err := DB.Exec("DELETE FROM anexos WHERE anexo_id = $1", id)
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
