package store

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

func HandleLojas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		lista, err := database.ListarLojas()
		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao listar lojas", http.StatusInternalServerError, err)
			return
		}
		if lista == nil {
			lista = []models.Loja{}
		}
		shared.RespondJSON(w, http.StatusOK, lista)
		return
	}

	if r.Method == "POST" {
		var l models.Loja
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if l.Nome == "" || l.LUC == "" || l.Segmento == "" {
			shared.RespondErrorJSON(w, "Os campos nome, luc e segmento são obrigatórios", http.StatusBadRequest, nil)
			return
		}

		id, err := database.CriarLoja(l)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				shared.RespondErrorJSON(w, "Já existe uma loja com este LUC", http.StatusConflict, err)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao criar loja", http.StatusInternalServerError, err)
			return
		}
		l.LojaID = id
		shared.RespondJSON(w, http.StatusCreated, l)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}

func HandleLoja(w http.ResponseWriter, r *http.Request) {
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
		l, err := database.BuscarLojaPorID(id)
		if err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Loja não encontrada", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao buscar loja", http.StatusInternalServerError, err)
			return
		}
		shared.RespondJSON(w, http.StatusOK, l)
		return
	}

	if r.Method == "PUT" {
		var l models.Loja
		if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}
		if l.Nome == "" || l.LUC == "" || l.Segmento == "" {
			shared.RespondErrorJSON(w, "Os campos nome, luc e segmento são obrigatórios", http.StatusBadRequest, nil)
			return
		}
		if err := database.AtualizarLoja(id, l); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Loja não encontrada", http.StatusNotFound, nil)
				return
			}
			if strings.Contains(err.Error(), "unique constraint") {
				shared.RespondErrorJSON(w, "Já existe uma loja com este LUC", http.StatusConflict, err)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao atualizar loja", http.StatusInternalServerError, err)
			return
		}
		l.LojaID = id
		shared.RespondJSON(w, http.StatusOK, l)
		return
	}

	if r.Method == "DELETE" {
		if err := database.DeletarLoja(id); err != nil {
			if err == sql.ErrNoRows {
				shared.RespondErrorJSON(w, "Loja não encontrada", http.StatusNotFound, nil)
				return
			}
			shared.RespondErrorJSON(w, "Erro ao excluir loja", http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	shared.RespondErrorJSON(w, "Method not allowed", http.StatusMethodNotAllowed, nil)
}
