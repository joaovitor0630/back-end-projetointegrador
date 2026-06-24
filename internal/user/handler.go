package user

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

func HandleUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		lista, err := database.ListarUsuarios()
		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao listar usuários", http.StatusInternalServerError, err)
			return
		}
		if lista == nil {
			lista = []models.Usuario{}
		}
		shared.RespondJSON(w, http.StatusOK, lista)
		return
	}

	if r.Method == "POST" {
		var u models.Usuario
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if u.Nome == "" || u.Departamento == "" || u.Cargo == "" {
			shared.RespondErrorJSON(w, "Os campos nome, departamento e cargo são obrigatórios", http.StatusBadRequest, nil)
			return
		}

		id, err := database.CriarUsuario(u)
		if err != nil {
			if strings.Contains(err.Error(), "foreign key constraint") {
				shared.RespondErrorJSON(w, "Perfil não encontrado", http.StatusBadRequest, err)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao criar usuário", http.StatusInternalServerError, err)
			return
		}
		u.UsuarioID = id
		shared.RespondJSON(w, http.StatusCreated, u)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}

func HandleUsuario(w http.ResponseWriter, r *http.Request) {
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
		u, err := database.BuscarUsuarioPorID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Usuário não encontrado", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao buscar usuário", http.StatusInternalServerError, err)
			return
		}
		shared.RespondJSON(w, http.StatusOK, u)
		return
	}

	if r.Method == "PUT" {
		var u models.Usuario
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if u.Nome == "" || u.Departamento == "" || u.Cargo == "" {
			shared.RespondErrorJSON(w, "Os campos nome, departamento e cargo são obrigatórios", http.StatusBadRequest, nil)
			return
		}

		if err := database.AtualizarUsuario(id, u); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Usuário não encontrado", http.StatusNotFound, nil)
				return
			}
			if strings.Contains(err.Error(), "foreign key constraint") {
				shared.RespondErrorJSON(w, "Perfil não encontrado", http.StatusBadRequest, err)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao atualizar usuário", http.StatusInternalServerError, err)
			return
		}
		u.UsuarioID = id
		shared.RespondJSON(w, http.StatusOK, u)
		return
	}

	if r.Method == "DELETE" {
		if err := database.DeletarUsuario(id); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Usuário não encontrado", http.StatusNotFound, nil)
				return
			}
			if strings.Contains(err.Error(), "não é possível excluir") {
				shared.RespondErrorJSON(w, err.Error(), http.StatusConflict, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao excluir usuário", http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}
