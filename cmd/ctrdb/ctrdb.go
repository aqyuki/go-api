package main

import (
	"context"
	"flag"
	"log/slog"

	"github.com/aqyuki/go-api/account"
	"github.com/aqyuki/go-api/infra/bbolt"
	"github.com/aqyuki/go-api/logging"
	"github.com/aqyuki/go-api/password"
	"github.com/aqyuki/go-api/prompt"
)

var (
	databasePathParm = flag.String("db", "", "path for database file")
)

func main() {
	logger := logging.NewLogger()

	if databasePathParm == nil {
		logger.Error("database path is not provided")
		return
	}
	databasePath := *databasePathParm
	logger.Info(
		"Starting ctrdb...",
		slog.String("databasePath", databasePath),
	)

	// initialize data
	id, err := prompt.StringFromPrompt("Input Account ID", false)
	if err != nil {
		logger.Error("failed to get account id", slog.Any("error", err))
		return
	}
	name, err := prompt.StringFromPrompt("Input Account Name", false)
	if err != nil {
		logger.Error("failed to get account name", slog.Any("error", err))
		return
	}
	pass, err := prompt.StringFromPrompt("Input Account Password", false)
	if err != nil {
		logger.Error("failed to get account password", slog.Any("error", err))
		return
	}
	bio, err := prompt.StringFromPrompt("Input Account Bio", true)
	if err != nil {
		logger.Error("failed to get account bio", slog.Any("error", err))
		return
	}

	encoder := &password.SHA256Encoder{}
	encoded, err := encoder.Encode(pass)
	if err != nil {
		logger.Error("failed to encode password", slog.Any("error", err))
		return
	}

	a := &account.Account{
		ID:           id,
		Name:         name,
		Bio:          bio,
		PasswordHash: encoded,
	}

	repo, err := bbolt.NewAccountRepository(databasePath)
	if err != nil {
		logger.Error("failed to initialize repository", slog.Any("error", err))
		return
	}
	defer repo.Close()

	err = repo.CreateAccount(context.Background(), a)
	if err != nil {
		logger.Error("failed to create account", slog.Any("error", err))
		return
	}
	logger.Info("account created", slog.Any("account", a))
}

func init() {
	flag.Parse()
}
