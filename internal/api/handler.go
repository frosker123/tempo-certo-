package api

import (
	"context"
	"net/http"

	"tempo/internal/container"
	"tempo/pkg/domains/empresas/transport"

	"github.com/go-chi/chi"
)

func Handler(ctx context.Context, dep *container.Dependency) http.Handler {
	r := chi.NewMux()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {})

	//criaçao de subi rotas da domains
	agendamentoHandler := transport.NewHTTPHandler(dep.Services.Agendamento)
	r.Mount("/agendamento", agendamentoHandler)

	return r
}
