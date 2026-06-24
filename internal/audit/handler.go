package audit

import (
	"encoding/json"
	"net/http"
	"strconv"
	"plataforma-flamboyant/internal/database"
	"plataforma-flamboyant/internal/models"
	"plataforma-flamboyant/internal/shared"
)

func HandleAuditoria(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tipo := r.URL.Query().Get("tipo")
		idStr := r.URL.Query().Get("id")
		
		if tipo == "" || idStr == "" {
			shared.RespondErrorJSON(w, "Faltando parâmetros 'tipo' e 'id'", http.StatusBadRequest, nil)
			return
		}
		
		id, err := strconv.Atoi(idStr)
		if err != nil {
			shared.RespondErrorJSON(w, "O 'id' deve ser numérico", http.StatusBadRequest, nil)
			return
		}
		
		logs, err := database.ListarAuditoriaPorEntidade(tipo, id)
		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao buscar logs", http.StatusInternalServerError, err)
			return
		}
		
		if logs == nil {
			logs = []models.Auditoria{} // evitar null JSON
		}
		
		shared.RespondJSON(w, http.StatusOK, logs)
		return
	}

	if r.Method == http.MethodPost {
		var a models.Auditoria
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			shared.RespondErrorJSON(w, "Dados inválidos", http.StatusBadRequest, err)
			return
		}

		if a.EntidadeTipo == "" || a.EntidadeID <= 0 || a.Acao == "" {
			shared.RespondErrorJSON(w, "Campos obrigatórios ausentes", http.StatusBadRequest, nil)
			return
		}

		id, err := database.CriarAuditoria(a)
		if err != nil {
			shared.RespondErrorJSON(w, "Erro ao criar log de auditoria", http.StatusInternalServerError, err)
			return
		}

		a.AuditoriaID = id
		shared.RespondJSON(w, http.StatusCreated, a)
		return
	}

	shared.RespondErrorJSON(w, "Método não permitido", http.StatusMethodNotAllowed, nil)
}
