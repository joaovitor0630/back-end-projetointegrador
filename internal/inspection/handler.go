package inspection

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

func HandleVistorias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		lista, err := database.ListarVistorias()
		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao listar vistorias", http.StatusInternalServerError, err)
			return
		}
		if lista == nil {
			lista = []models.Vistoria{}
		}
		shared.RespondJSON(w, http.StatusOK, lista)
		return
	}

	if r.Method == "POST" {
		var v models.Vistoria
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if v.LojaID <= 0 || v.UsuarioID <= 0 || v.Descricao == "" {
			shared.RespondErrorJSON(w, "Os campos loja_id, usuario_id e descricao são obrigatórios", http.StatusBadRequest, nil)
			return
		}
		if !v.AreaResponsavel.IsValid() {
			shared.RespondErrorJSON(w, "Área responsável inválida", http.StatusBadRequest, nil)
			return
		}
		if !v.Categoria.IsValid() {
			shared.RespondErrorJSON(w, "Categoria inválida", http.StatusBadRequest, nil)
			return
		}

		id, err := database.CriarVistoria(v)
		if err != nil {
			if strings.Contains(err.Error(), "foreign key constraint") {
				shared.RespondErrorJSON(w, "Loja ou Usuário não encontrado", http.StatusBadRequest, err)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao criar vistoria", http.StatusInternalServerError, err)
			return
		}
		v.VistoriaID = id
		shared.RespondJSON(w, http.StatusCreated, v)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}

func HandleVistoria(w http.ResponseWriter, r *http.Request) {
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
		v, err := database.BuscarVistoriaPorID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Vistoria não encontrada", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao buscar vistoria", http.StatusInternalServerError, err)
			return
		}
		shared.RespondJSON(w, http.StatusOK, v)
		return
	}

	if r.Method == "PUT" {
		var v models.Vistoria
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if v.LojaID <= 0 || v.UsuarioID <= 0 || v.Descricao == "" {
			shared.RespondErrorJSON(w, "Os campos loja_id, usuario_id e descricao são obrigatórios", http.StatusBadRequest, nil)
			return
		}
		if !v.AreaResponsavel.IsValid() {
			shared.RespondErrorJSON(w, "Área responsável inválida", http.StatusBadRequest, nil)
			return
		}
		if !v.Categoria.IsValid() {
			shared.RespondErrorJSON(w, "Categoria inválida", http.StatusBadRequest, nil)
			return
		}

		if err := database.AtualizarVistoria(id, v); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Vistoria não encontrada", http.StatusNotFound, nil)
				return
			}
			if strings.Contains(err.Error(), "foreign key constraint") {
				shared.RespondErrorJSON(w, "Loja ou Usuário não encontrado", http.StatusBadRequest, err)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao atualizar vistoria", http.StatusInternalServerError, err)
			return
		}
		v.VistoriaID = id
		shared.RespondJSON(w, http.StatusOK, v)
		return
	}

	if r.Method == "DELETE" {
		if err := database.DeletarVistoria(id); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Vistoria não encontrada", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao excluir vistoria", http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}