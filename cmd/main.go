package main

import "library/internal/logger"

func main() {
	logger, logFileClose, err := logger.NewLogger("INFO")
	if err != nil {
		panic(err)
	}
	defer logFileClose()
	logger.Info("app started")
}
