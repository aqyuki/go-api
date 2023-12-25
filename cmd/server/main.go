package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/aqyuki/go-api/account"
	"github.com/aqyuki/go-api/infra/bbolt"
	"github.com/aqyuki/go-api/logging"
	"github.com/aqyuki/go-api/password"
	"github.com/aqyuki/go-api/server"
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
	defer repo.Close()

	service := account.NewAccountApp(repo, &password.SHA256Encoder{})
	server := server.NewAccountServer([]byte(secret), logger, service)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		logger.Info("Shutting down server...", slog.String("reason", ctx.Err().Error()))
		cancelCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		server.Shutdown(cancelCtx)
		logger.Info("Server shutdown")
	}()
	server.Start(":8080")
}
