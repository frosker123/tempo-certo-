package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"tempo/internal/api"
	"tempo/internal/container"
	"tempo/pkg/core/env"
	apiServer "tempo/pkg/core/http"
	"tempo/pkg/core/types"
)

func main() {

	// criaçao  meu context
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.ContextKey(types.StartedAt), time.Now())
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	// carrega as variaveis de ambiente
	err := env.LoadEnv(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//  carrega as dependencias do container
	ctx, dep, err := container.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// cria o servidor e roda a api
	apiServer.Run(
		ctx,
		fmt.Sprintf("%v", env.API_PORT),
		api.Handler(ctx, dep),
	)

}
