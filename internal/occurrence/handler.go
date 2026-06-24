package occurrence

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"plataforma-flamboyant/internal/database"
	"plataforma-flamboyant/internal/models"
	"plataforma-flamboyant/internal/shared"
	"strconv"
	"strings"
)

func HandleOcorrencias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		lista, err := database.ListarOcorrencias()
		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao listar ocorrências", http.StatusInternalServerError, err)
			return
		}
		if lista == nil {
			lista = []models.Ocorrencia{}
		}
		shared.RespondJSON(w, http.StatusOK, lista)
		return
	}

	if r.Method == "POST" {
		var o models.Ocorrencia
		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if o.LojaID <= 0 || o.Assunto == "" || o.Descricao == "" {
			shared.RespondErrorJSON(w, "Os campos loja_id, assunto e descricao são obrigatórios", http.StatusBadRequest, nil)
			return
		}
		if !o.AreaResponsavel.IsValid() {
			shared.RespondErrorJSON(w, "Área responsável inválida", http.StatusBadRequest, nil)
			return
		}
		if !o.Categoria.IsValid() {
			shared.RespondErrorJSON(w, "Categoria inválida", http.StatusBadRequest, nil)
			return
		}

		id, err := database.CriarOcorrencia(o)
		if err != nil {
			if strings.Contains(err.Error(), "foreign key constraint") {
				shared.RespondErrorJSON(w, "Loja não encontrada", http.StatusBadRequest, err)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao criar ocorrência", http.StatusInternalServerError, err)
			return
		}
		o.OcorrenciaID = id
		shared.RespondJSON(w, http.StatusCreated, o)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}

func HandleOcorrencia(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "ID missing", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		o, err := database.BuscarOcorrenciaPorID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Ocorrência não encontrada", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao buscar ocorrência", http.StatusInternalServerError, err)
			return
		}
		shared.RespondJSON(w, http.StatusOK, o)
		return
	}

	if r.Method == "PUT" {
		var o models.Ocorrencia
		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if o.LojaID <= 0 || o.Assunto == "" || o.Descricao == "" {
			shared.RespondErrorJSON(w, "Os campos loja_id, assunto e descricao são obrigatórios", http.StatusBadRequest, nil)
			return
		}
		if !o.AreaResponsavel.IsValid() {
			shared.RespondErrorJSON(w, "Área responsável inválida", http.StatusBadRequest, nil)
			return
		}
		if !o.Categoria.IsValid() {
			shared.RespondErrorJSON(w, "Categoria inválida", http.StatusBadRequest, nil)
			return
		}

		if err := database.AtualizarOcorrencia(id, o); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Ocorrência não encontrada", http.StatusNotFound, nil)
				return
			}
			if strings.Contains(err.Error(), "foreign key constraint") {
				shared.RespondErrorJSON(w, "Loja não encontrada", http.StatusBadRequest, err)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao atualizar ocorrência", http.StatusInternalServerError, err)
			return
		}
		o.OcorrenciaID = id
		shared.RespondJSON(w, http.StatusOK, o)
		return
	}

	if r.Method == "DELETE" {
		if err := database.DeletarOcorrencia(id); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Ocorrência não encontrada", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao excluir ocorrência", http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}