package profile

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

func HandlePerfis(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		lista, err := database.ListarPerfis()
		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao listar perfis", http.StatusInternalServerError, err)
			return
		}
		if lista == nil {
			lista = []models.Perfil{}
		}
		shared.RespondJSON(w, http.StatusOK, lista)
		return
	}

	if r.Method == "POST" {
		var p models.Perfil
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if p.TipoPerfil == "" || p.Descricao == "" {
			shared.RespondErrorJSON(w, "Os campos tipo_perfil e descricao são obrigatórios", http.StatusBadRequest, nil)
			return
		}

		id, err := database.CriarPerfil(p)
		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao criar perfil", http.StatusInternalServerError, err)
			return
		}
		p.PerfilID = id
		shared.RespondJSON(w, http.StatusCreated, p)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}

func HandlePerfil(w http.ResponseWriter, r *http.Request) {
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
		p, err := database.BuscarPerfilPorID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Perfil não encontrado", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao buscar perfil", http.StatusInternalServerError, err)
			return
		}
		shared.RespondJSON(w, http.StatusOK, p)
		return
	}

	if r.Method == "PUT" {
		var p models.Perfil
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if p.TipoPerfil == "" || p.Descricao == "" {
			shared.RespondErrorJSON(w, "Os campos tipo_perfil e descricao são obrigatórios", http.StatusBadRequest, nil)
			return
		}

		if err := database.AtualizarPerfil(id, p); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Perfil não encontrado", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao atualizar perfil", http.StatusInternalServerError, err)
			return
		}
		p.PerfilID = id
		shared.RespondJSON(w, http.StatusOK, p)
		return
	}

	if r.Method == "DELETE" {
		if err := database.DeletarPerfil(id); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Perfil não encontrado", http.StatusNotFound, nil)
				return
			}
			if strings.Contains(err.Error(), "não é possível excluir") {
				shared.RespondErrorJSON(w, err.Error(), http.StatusConflict, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao excluir perfil", http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}
