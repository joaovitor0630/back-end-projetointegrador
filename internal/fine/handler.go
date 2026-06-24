package fine

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

func HandleMultas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		lista, err := database.ListarMultas()
		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao listar multas", http.StatusInternalServerError, err)
			return
		}
		if lista == nil {
			lista = []models.Multa{}
		}
		shared.RespondJSON(w, http.StatusOK, lista)
		return
	}

	if r.Method == "POST" {
		var m models.Multa
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if m.OcorrenciaID <= 0 || m.Assunto == "" {
			shared.RespondErrorJSON(w, "Os campos ocorrencia_id e assunto são obrigatórios", http.StatusBadRequest, nil)
			return
		}
		if !m.Categoria.IsValid() {
			shared.RespondErrorJSON(w, "Categoria inválida", http.StatusBadRequest, nil)
			return
		}

		if m.Status == "" {
			m.Status = models.StatusPendente
		} else if !m.Status.IsValid() {
			shared.RespondErrorJSON(w, "Status inválido", http.StatusBadRequest, nil)
			return
		}

		id, err := database.CriarMulta(m)
		if err != nil {
			if strings.Contains(err.Error(), "foreign key constraint") {
				shared.RespondErrorJSON(w, "Ocorrência não encontrada", http.StatusBadRequest, err)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao criar multa", http.StatusInternalServerError, err)
			return
		}
		m.MultaID = id
		shared.RespondJSON(w, http.StatusCreated, m)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}

func HandleMulta(w http.ResponseWriter, r *http.Request) {
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
		m, err := database.BuscarMultaPorID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Multa não encontrada", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao buscar multa", http.StatusInternalServerError, err)
			return
		}
		shared.RespondJSON(w, http.StatusOK, m)
		return
	}

	if r.Method == "PUT" {
		var m models.Multa
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if m.OcorrenciaID <= 0 || m.Assunto == "" {
			shared.RespondErrorJSON(w, "Os campos ocorrencia_id e assunto são obrigatórios", http.StatusBadRequest, nil)
			return
		}
		if !m.Categoria.IsValid() {
			shared.RespondErrorJSON(w, "Categoria inválida", http.StatusBadRequest, nil)
			return
		}

		if m.Status == "" {
			m.Status = models.StatusPendente
		} else if !m.Status.IsValid() {
			shared.RespondErrorJSON(w, "Status inválido", http.StatusBadRequest, nil)
			return
		}

		if err := database.AtualizarMulta(id, m); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Multa não encontrada", http.StatusNotFound, nil)
				return
			}
			if strings.Contains(err.Error(), "foreign key constraint") {
				shared.RespondErrorJSON(w, "Ocorrência não encontrada", http.StatusBadRequest, err)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao atualizar multa", http.StatusInternalServerError, err)
			return
		}
		m.MultaID = id
		shared.RespondJSON(w, http.StatusOK, m)
		return
	}

	if r.Method == "DELETE" {
		if err := database.DeletarMulta(id); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Multa não encontrada", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao excluir multa", http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}