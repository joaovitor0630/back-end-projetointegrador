package attachment

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"plataforma-flamboyant/internal/database"
	"plataforma-flamboyant/internal/models"
	"plataforma-flamboyant/internal/shared"
	"strconv"
	"strings"
)

func HandleAnexos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		q := r.URL.Query()
		var lista []models.Anexo
		var err error

		if ocorrenciaID := q.Get("ocorrencia_id"); ocorrenciaID != "" {
			id, convErr := strconv.Atoi(ocorrenciaID)
			if convErr != nil {
				shared.RespondErrorJSON(w, "ocorrencia_id inválido", http.StatusBadRequest, nil)
				return
			}
			lista, err = database.ListarAnexosPorOcorrencia(id)
		} else if multaID := q.Get("multa_id"); multaID != "" {
			id, convErr := strconv.Atoi(multaID)
			if convErr != nil {
				shared.RespondErrorJSON(w, "multa_id inválido", http.StatusBadRequest, nil)
				return
			}
			lista, err = database.ListarAnexosPorMulta(id)
		} else if vistoriaID := q.Get("vistoria_id"); vistoriaID != "" {
			id, convErr := strconv.Atoi(vistoriaID)
			if convErr != nil {
				shared.RespondErrorJSON(w, "vistoria_id inválido", http.StatusBadRequest, nil)
				return
			}
			lista, err = database.ListarAnexosPorVistoria(id)
		} else {
			shared.RespondErrorJSON(w, "É necessário passar ocorrencia_id, multa_id ou vistoria_id", http.StatusBadRequest, nil)
			return
		}

		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao listar anexos", http.StatusInternalServerError, err)
			return
		}
		if lista == nil {
			lista = []models.Anexo{}
		}
		for i := range lista {
			lista[i].Conteudo = nil
		}
		shared.RespondJSON(w, http.StatusOK, lista)
		return
	}

	if r.Method == "POST" {
		var a models.Anexo
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if a.NomeArquivo == "" {
			shared.RespondErrorJSON(w, "nome_arquivo é obrigatório", http.StatusBadRequest, nil)
			return
		}

		id, err := database.CriarAnexo(a)
		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao criar anexo", http.StatusInternalServerError, err)
			return
		}
		a.AnexoID = id
		a.Conteudo = nil // don't echo binary back
		shared.RespondJSON(w, http.StatusCreated, a)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}

func HandleAnexo(w http.ResponseWriter, r *http.Request) {
	// Path: /api/anexos/{id} or /api/anexos/{id}/download
	parts := strings.Split(strings.TrimRight(r.URL.Path, "/"), "/")
	// parts: ["", "api", "anexos", "{id}"] or ["", "api", "anexos", "{id}", "download"]
	if len(parts) < 4 {
		shared.RespondErrorJSON(w, "ID missing", http.StatusBadRequest, nil)
		return
	}

	isDownload := len(parts) >= 5 && parts[4] == "download"
	id, err := strconv.Atoi(parts[3])
	if err != nil {
		shared.RespondErrorJSON(w, "Invalid ID", http.StatusBadRequest, nil)
		return
	}

	// GET /api/anexos/:id/download — serve raw binary with correct Content-Type
	if r.Method == "GET" && isDownload {
		a, err := database.BuscarAnexoPorID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Anexo não encontrado", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao buscar anexo", http.StatusInternalServerError, err)
			return
		}
		mime := a.TipoMime
		if mime == "" {
			mime = "application/octet-stream"
		}
		safeFilename := strings.ReplaceAll(a.NomeArquivo, "\"", "\\\"")
		safeFilename = strings.ReplaceAll(safeFilename, "\r", "")
		safeFilename = strings.ReplaceAll(safeFilename, "\n", "")
		w.Header().Set("Content-Type", mime)
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", safeFilename))
		w.Header().Set("Content-Length", strconv.Itoa(len(a.Conteudo)))
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(a.Conteudo)
		return
	}

	// GET /api/anexos/:id — JSON metadata (without binary content)
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		a, err := database.BuscarAnexoPorID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Anexo não encontrado", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao buscar anexo", http.StatusInternalServerError, err)
			return
		}
		a.Conteudo = nil // omit binary from JSON
		shared.RespondJSON(w, http.StatusOK, a)
		return
	}

	if r.Method == "PUT" {
		w.Header().Set("Content-Type", "application/json")
		var a models.Anexo
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}
		if a.NomeArquivo == "" {
			shared.RespondErrorJSON(w, "nome_arquivo é obrigatório", http.StatusBadRequest, nil)
			return
		}
		if err := database.AtualizarAnexo(id, a); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Anexo não encontrado", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao atualizar anexo", http.StatusInternalServerError, err)
			return
		}
		a.AnexoID = id
		a.Conteudo = nil // don't echo binary content on update
		shared.RespondJSON(w, http.StatusOK, a)
		return
	}

	if r.Method == "DELETE" {
		if err := database.DeletarAnexo(id); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Anexo não encontrado", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao excluir anexo", http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}
