/*
Package server is the package responsible for http server
*/
package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tempo/pkg/core/env"
	"tempo/pkg/core/types"
)

// PrefixConfig is the environment prefix ex: HOST_PORT will be added API_ as prefix.
const PrefixConfig = "API_"

// Run the server http.
func Run(ctx context.Context, port string, handler http.Handler) error {
	env.LoadEnv(ctx)
	var GracefulDuration time.Duration

	server := &http.Server{Addr: fmt.Sprintf(":%v", env.API_PORT), Handler: handler}

	// Server run context.
	serverCtx, serverStopCtx := context.WithCancel(ctx)

	// Listen for syscall signals for process to interrupt/quit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, GracefulDuration)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Println("HTTP server graceful shutdown timeout")
			}
		}()

		// Trigger graceful shutdown
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Println("HTTP server graceful shutdown error:", err)
		}
		serverStopCtx()
	}()

	startingMessage(ctx, fmt.Sprintf(":%v", env.API_PORT))

	if err := server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		log.Println("HTTP server graceful shutdown")
	} else {
		log.Println("HTTP server error:", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	return nil
}

func startingMessage(ctx context.Context, where string) {
	t, ok := ctx.Value(types.ContextKey(types.StartedAt)).(time.Time)
	if !ok {
		return
	}

	fmt.Printf("HTTP server started at %v on %v", t.Format(time.Kitchen), where)
}
