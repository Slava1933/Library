package main

import (
	"context"
	"library/internal/logger"
	"library/internal/repository"
)

func main() {
	logger, logFileClose, err := logger.NewLogger("INFO")
	if err != nil {
		panic(err)
	}
	defer logFileClose()
	logger.Info("app started")

	ctx := context.Background()

	pool, err := repository.Connect(ctx)
	if err != nil {
		logger.Fatal("Не удалось подключиться к БД")
	}
	defer pool.Close()

}
