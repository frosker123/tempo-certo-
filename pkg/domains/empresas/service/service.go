package service

import (
	"context"
	"tempo/pkg/domains/empresas/model"
	"tempo/pkg/domains/empresas/repository"

	"github.com/pkg/errors"
)

type ServiceI interface {
	Upsert(context.Context, *model.Agenda) error
	CreateListaHorario(context.Context, []model.HorarioAgendamento) error
	ListAgendamentos(context.Context) ([]model.ListAgendados, error)
	ListAgendaDisponibilidade(context.Context) ([]model.HorarioAgendamento, error)
	FindByCnpj(context.Context, string) ([]model.ListAgendados, error)
}

type Service struct {
	repository repository.RepositoryI
}

func NewService(repositorio repository.RepositoryI) (*Service, error) {
	if repositorio == nil {
		return nil, errors.New("repository is nil")
	}
	return &Service{
		repository: repositorio,
	}, nil
}

func (s *Service) Upsert(ctx context.Context, upsert *model.Agenda) error {
	return s.repository.Upsert(ctx, upsert)
}

func (s *Service) ListAgendamentos(ctx context.Context) ([]model.ListAgendados, error) {
	return s.repository.ListAgendamentos(ctx)
}

func (s *Service) ListAgendaDisponibilidade(ctx context.Context) ([]model.HorarioAgendamento, error) {
	return s.repository.ListAgendaDisponibilidade(ctx)
}

func (s *Service) FindByCnpj(ctx context.Context, cnpj string) ([]model.ListAgendados, error) {
	return s.repository.FindByCnpj(ctx, cnpj)
}

func (s *Service) CreateListaHorario(ctx context.Context, lista []model.HorarioAgendamento) error {
	return s.repository.CreateListaHorario(ctx, lista)
}
