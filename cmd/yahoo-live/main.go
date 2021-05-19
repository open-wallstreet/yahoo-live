package main

import (
	"os"

	"github.com/open-wallstreet/yahoo-live/internal/yahoo"
	"github.com/open-wallstreet/yahoo-live/proto"
	"go.uber.org/zap"
)

func createLogger() *zap.SugaredLogger {
	appEnv := os.Getenv("APP_ENV")
	var logger *zap.Logger
	if appEnv == "PRODUCTION" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	sugar := logger.Sugar()
	return sugar
}

func on_msg(message *proto.Yaticker) {
	println(message.String())
}

func main() {
	logger := createLogger()
	yahoo.NewYahooWebsocket(logger, []string{"AMZN", "AAPL", "TSLA", "A", "AA"}, on_msg)
}
