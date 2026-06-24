package main

import (
	"log"
	"net/http"
	"plataforma-flamboyant/internal/audit"
	"plataforma-flamboyant/internal/attachment"
	"plataforma-flamboyant/internal/database"
	"plataforma-flamboyant/internal/fine"
	"plataforma-flamboyant/internal/inspection"
	"plataforma-flamboyant/internal/middleware"
	"plataforma-flamboyant/internal/occurrence"
	"plataforma-flamboyant/internal/profile"
	"plataforma-flamboyant/internal/store"
	"plataforma-flamboyant/internal/user"
)

func main() {
	// Connect to database
	database.ConnectToDB()

	// Use default mux
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/api/lojas", store.HandleLojas)
	mux.HandleFunc("/api/lojas/", store.HandleLoja)

	mux.HandleFunc("/api/ocorrencias", occurrence.HandleOcorrencias)
	mux.HandleFunc("/api/ocorrencias/", occurrence.HandleOcorrencia)

	mux.HandleFunc("/api/multas", fine.HandleMultas)
	mux.HandleFunc("/api/multas/", fine.HandleMulta)

	mux.HandleFunc("/api/vistorias", inspection.HandleVistorias)
	mux.HandleFunc("/api/vistorias/", inspection.HandleVistoria)

	mux.HandleFunc("/api/usuarios", user.HandleUsuarios)
	mux.HandleFunc("/api/usuarios/", user.HandleUsuario)

	mux.HandleFunc("/api/perfis", profile.HandlePerfis)
	mux.HandleFunc("/api/perfis/", profile.HandlePerfil)

	mux.HandleFunc("/api/anexos", attachment.HandleAnexos)
	mux.HandleFunc("/api/anexos/", attachment.HandleAnexo)
	
	mux.HandleFunc("/api/auditoria", audit.HandleAuditoria)

	// Status endpoint for Dashboard health check
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"mensagem": "API Online"}`))
	})

	// Apply middlewares: Logger -> SecurityHeaders -> MaxBodySize -> CORS -> Mux
	handler := middleware.Wrap(mux.ServeHTTP, middleware.CORS, middleware.MaxBodySize, middleware.SecurityHeaders, middleware.Logger)

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
