package transport

import (
	"context"
	"encoding/json"
	"log"
	stdHTTP "net/http"

	"tempo/pkg/domains/empresas/model"
	"tempo/pkg/domains/empresas/service"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/transport/http"
)

func NewHTTPHandler(svc service.ServiceI) stdHTTP.Handler {
	options := []http.ServerOption{
		http.ServerErrorEncoder(ResponseErrors),
	}

	upsert := http.NewServer(
		service.Upsert(svc),
		decodeUpsertRequest,
		codeHTTP{200}.encodeResponse,
		options...,
	)

	listAgendamentos := http.NewServer(
		service.ListAgendamentos(svc),
		http.NopRequestDecoder,
		codeHTTP{200}.encodeResponse,
		options...,
	)

	findByCnpj := http.NewServer(
		service.FindByCnpj(svc),
		decodeFindRequest,
		codeHTTP{200}.encodeResponse,
		options...,
	)

	listAgendaDisponibilidade := http.NewServer(
		service.ListAgendaDisponibilidade(svc),
		http.NopRequestDecoder,
		codeHTTP{200}.encodeResponse,
		options...,
	)
	createlista := http.NewServer(
		service.CreateListaHorario(svc),
		http.NopRequestDecoder,
		codeHTTP{200}.encodeResponse,
		options...,
	)

	r := chi.NewRouter()

	r.Post("/agendas", upsert.ServeHTTP)
	r.Post("/lista", createlista.ServeHTTP)
	r.Get("/agendas/", listAgendamentos.ServeHTTP)
	r.Get("/agendas/{cnpj}", findByCnpj.ServeHTTP)
	r.Get("/agendas/disponibilidade", listAgendaDisponibilidade.ServeHTTP)

	return r
}

func decodeUpsertRequest(_ context.Context, r *stdHTTP.Request) (interface{}, error) {
	req := model.UpsertRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeFindRequest(_ context.Context, r *stdHTTP.Request) (interface{}, error) {
	return model.ListAgendados{
		Cnpj: chi.URLParam(r, "cnpj"),
	}, nil
}

type codeHTTP struct {
	int
}
type Errors struct {
	Erros string `json:"Erros"`
}

func (c codeHTTP) encodeResponse(_ context.Context, w stdHTTP.ResponseWriter, input interface{}) error {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(c.int)
	return json.NewEncoder(w).Encode(input)
}

func ResponseErrors(_ctx context.Context, err error, w stdHTTP.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(stdHTTP.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(Errors{Erros: err.Error()}); err != nil {
		log.Fatal("erro no encode")
		return
	}

}
