package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tempo/pkg/core/env"
	"tempo/pkg/domains/empresas/model"
)

type RepositoryI interface {
	// Queries is a "Readeble" interface responsible to read data from source
	Querier

	// Execer is a "Writable" interface responsible for write data into source
	Execer
}

type Querier interface {
	ListAgendamentos(ctx context.Context) ([]model.ListAgendados, error)
	ListAgendaDisponibilidade(ctx context.Context) ([]model.HorarioAgendamento, error)
	FindByCnpj(ctx context.Context, cnpj string) ([]model.ListAgendados, error)
}

type Execer interface {
	Upsert(context.Context, *model.Agenda) error
	CreateListaHorario(context.Context, []model.HorarioAgendamento) error
}

type RepositoryMemory struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *RepositoryMemory {
	return &RepositoryMemory{db}
}

func (r *RepositoryMemory) CreateListaHorario(ctx context.Context, lista []model.HorarioAgendamento) error {
	err := r.createTableLista()
	if err != nil {
		return err
	}

	inserir, err := r.db.Prepare("insert into disponibilidade(horario_agendado,horario_termino,disponivel)values($1, $2, $3)")
	if err != nil {
		return err
	}

	for _, item := range lista {
		if _, err = inserir.Exec(item.Inicio, item.Fim, item.Disponivel); err != nil {
			return err
		}
	}
	return nil
}

func (r *RepositoryMemory) Upsert(ctx context.Context, agendamento *model.Agenda) error {
	err := r.createTable()
	if err != nil {
		return err
	}
	update, err := r.db.Prepare("UPDATE agendamento SET horario_agendado = $1  WHERE cnpj = $2")
	if err != nil {
		return err
	}
	result, err := update.Exec(agendamento.Horario, agendamento.Cnpj)
	if err != nil {
		return err
	}

	affect, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affect == 0 {
		insert, err := r.db.Prepare("insert into agendamento(fantasia,cnpj,horario_agendado)values($1, $2, $3)")
		if err != nil {
			return err
		}

		_, err = insert.Exec(agendamento.Empresa.Fantasia, agendamento.Empresa.Cnpj, agendamento.Empresa.Horario)
		if err != nil {
			return err
		}
	}

	updates, err := r.db.Prepare("UPDATE disponibilidade SET disponivel = $1 WHERE horario_agendado = $2")
	if err != nil {
		return err
	}
	_, err = updates.Exec(false, agendamento.Horario)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryMemory) findByCnpj(ctx context.Context, cnpj string) ([]model.ListAgendados, error) {
	var agendados []model.ListAgendados
	err := r.createTable()
	if err != nil {
		return []model.ListAgendados{}, err
	}

	row, err := r.db.Query("SELECT cnpj, fantasia, horario_agendado FROM agendamento WHERE cnpj = $1", cnpj)
	if err != nil {
		return []model.ListAgendados{}, err
	}
	for row.Next() {
		var findRazao model.ListAgendados
		if err = row.Scan(&findRazao.Cnpj, &findRazao.Fantasia, &findRazao.Inicio); err != nil {
			return nil, err
		}

		findRazao, err = completeRazao(ctx, findRazao)
		if err != nil {
			return nil, err
		}

		agendados = append(agendados, findRazao)
	}

	return agendados, nil
}
func (r *RepositoryMemory) FindByCnpj(ctx context.Context, cnpj string) ([]model.ListAgendados, error) {
	return r.findByCnpj(ctx, cnpj)
}

func completeRazao(ctx context.Context, findRazao model.ListAgendados) (model.ListAgendados, error) {

	var (
		url = fmt.Sprintf("%s%s", env.ENDOPOINT_API, findRazao.Cnpj)
	)

	if findRazao.Fantasia == "" {
		req, err := http.Get(url)
		if err != nil {
			return findRazao, err
		}

		response, err := io.ReadAll(req.Body)
		if err != nil {
			return findRazao, err
		}
		err = json.Unmarshal(response, &findRazao)
		if err != nil {
			return findRazao, err
		}

	}
	return findRazao, nil
}

func (r *RepositoryMemory) ListAgendamentos(ctx context.Context) ([]model.ListAgendados, error) {
	var agendados []model.ListAgendados

	rows, err := r.db.Query("SELECT  cnpj, fantasia, horario_agendado FROM agendamento")
	if err != nil {
		return []model.ListAgendados{}, err
	}

	for rows.Next() {
		var findRazao model.ListAgendados

		err = rows.Scan(&findRazao.Cnpj, &findRazao.Fantasia, &findRazao.Inicio)
		if err != nil {
			return nil, err
		}

		findRazao, err = completeRazao(ctx, findRazao)
		if err != nil {
			return nil, err
		}

		agendados = append(agendados, findRazao)
	}
	return agendados, nil
}

func (r *RepositoryMemory) ListAgendaDisponibilidade(ctx context.Context) ([]model.HorarioAgendamento, error) {
	var disponibilidades []model.HorarioAgendamento

	rows, err := r.db.Query("SELECT horario_agendado,horario_termino, disponivel FROM disponibilidade")
	if err != nil {
		return []model.HorarioAgendamento{}, err
	}

	for rows.Next() {
		var disponibilidade model.HorarioAgendamento

		err = rows.Scan(&disponibilidade.Inicio, &disponibilidade.Fim, &disponibilidade.Disponivel)
		if err != nil {
			return []model.HorarioAgendamento{}, nil
		}
		disponibilidades = append(disponibilidades, disponibilidade)
	}
	return disponibilidades, nil
}

func (r *RepositoryMemory) createTable() error {

	create, err := r.db.Prepare("Create table if not exists agendamento(id serial primary key, fantasia varchar(200),cnpj varchar(14), horario_agendado varchar(10))")
	if err != nil {
		return err
	}

	_, err = create.Exec()
	if err != nil {
		return err
	}
	return nil
}
func (r *RepositoryMemory) createTableLista() error {

	create, err := r.db.Prepare("Create table if not exists disponibilidade(id serial primary key, horario_agendado varchar(10), horario_termino varchar(10), disponivel boolean)")
	if err != nil {
		return err
	}

	_, err = create.Exec()
	if err != nil {
		return err
	}
	return nil
}
