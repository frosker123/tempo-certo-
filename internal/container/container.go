package container

import (
	"context"
	"database/sql"
	"log"
	"tempo/pkg/core/env"
	"tempo/pkg/domains/empresas/repository"
	"tempo/pkg/domains/empresas/service"

	_ "github.com/lib/pq"
)

type components struct {
	db *sql.DB
}

type Services struct {
	Agendamento service.ServiceI
}

type Dependency struct {
	Components components
	Services   Services
}

func New(ctx context.Context) (context.Context, *Dependency, error) {
	cmp, err := setupComponents(ctx)
	if err != nil {
		return nil, nil, err
	}

	agedamentoService, err := service.NewService(
		repository.NewRepository(cmp.db),
	)
	if err != nil {
		return nil, nil, err
	}

	srv := Services{
		agedamentoService,
	}

	dep := Dependency{
		Components: *cmp,
		Services:   srv,
	}

	return ctx, &dep, err
}

func setupComponents(ctx context.Context) (*components, error) {

	db, err := sql.Open(env.DATABASE_DRIVE, env.CONECT)
	if err != nil {
		log.Fatal("erro ao conetar com o banco de dados ")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("erro no ping ")
		return nil, err
	}

	return &components{
		db: db,
	}, nil
}
