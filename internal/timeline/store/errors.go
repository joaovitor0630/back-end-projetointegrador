package store

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondError envia uma mensagem de erro genérica ao cliente e loga o erro real.
// Isso evita expor detalhes internos do banco de dados ou do sistema ao usuário.
func RespondError(w http.ResponseWriter, publicMsg string, statusCode int, internalErr error) {
	if internalErr != nil {
		log.Printf("[ERRO %d] %s: %v", statusCode, publicMsg, internalErr)
	}
	http.Error(w, publicMsg, statusCode)
}

// RespondJSON envia uma resposta JSON padronizada.
func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Erro ao encodar JSON de resposta: %v", err)
		}
	}
}

// RespondErrorJSON envia uma mensagem de erro em formato JSON.
func RespondErrorJSON(w http.ResponseWriter, publicMsg string, statusCode int, internalErr error) {
	if internalErr != nil {
		log.Printf("[ERRO %d] %s: %v", statusCode, publicMsg, internalErr)
	}
	RespondJSON(w, statusCode, map[string]string{"erro": publicMsg})
}

// IsRequiredEmpty verifica se algum dos campos num map está vazio.
// Retorna o nome do primeiro campo vazio encontrado, e um booleano (true se algum vazio).
func IsRequiredEmpty(fields map[string]string) (string, bool) {
	for name, val := range fields {
		if val == "" {
			return name, true
		}
	}
	return "", false
}