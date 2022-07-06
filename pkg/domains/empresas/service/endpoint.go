package service

import (
	"context"
	"tempo/pkg/domains/empresas/model"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

func Upsert(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.UpsertRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		response, err := model.NewAgendamento(req.Fantasia, req.Cnpj, req.Horario)
		if err != nil {
			return nil, err
		}

		if err := svc.Upsert(ctx, &response); err != nil {
			return nil, err
		}
		return response, nil
	}
}
func CreateListaHorario(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {

		response, err := model.NewListaDisponibilidade()
		if err != nil {
			return nil, err
		}

		if err := svc.CreateListaHorario(ctx, response); err != nil {
			return nil, err
		}
		return response, nil
	}
}

func ListAgendamentos(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.ListAgendamentos(ctx)
	}
}

func FindByCnpj(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(model.ListAgendados)
		if !ok {
			return nil, errors.New("invalid request")
		}

		response, err := svc.FindByCnpj(ctx, req.Cnpj)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}

func ListAgendaDisponibilidade(svc ServiceI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.ListAgendaDisponibilidade(ctx)
	}
}
