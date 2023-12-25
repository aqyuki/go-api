package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/aqyuki/jwt-demo/account"
	"github.com/aqyuki/jwt-demo/infra/bbolt"
	"github.com/aqyuki/jwt-demo/logging"
	"github.com/aqyuki/jwt-demo/password"
	"github.com/aqyuki/jwt-demo/server"
)

const (
	secret = "secret"
)

func main() {

	// initialize required modules
	logger := logging.NewLogger()
	logger.Info("Starting server...")

	// initialize server
	repo, err := bbolt.NewAccountRepository("test.db")
	if err != nil {
		logger.Error("failed to initialize repository", slog.Any("error", err))
		os.Exit(1)
	}
	service := account.NewAccountApp(repo, &password.SHA256Encoder{})
	server := server.NewAccountServer([]byte(secret), logger, service)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		logger.Info("Shutting down server...")
		cancelCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		server.Shutdown(cancelCtx)
		logger.Info("Server shutdown")
	}()
	server.Start(":8080")
}