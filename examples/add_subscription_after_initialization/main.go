package main

import (
	"fmt"
	"os"
	"time"

	"github.com/open-wallstreet/yahoo-live/pkg/yahoo"
	"github.com/open-wallstreet/yahoo-live/pkg/yahoo/proto"
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
	println(fmt.Sprintf("%s: %s", time.Unix(message.Time/1000, 0).String(), message.String()))
}

func main() {
	logger := createLogger()
	con, err := yahoo.NewWebsocket(logger, []string{"KIND-SDB.ST"})
	if err != nil {
		panic(err)
	}
	con.AddMessageHandler(on_msg)


    time.Sleep(2*time.Second)
    con.AddSubscription([]string{"TSLA"})
    time.Sleep(2*time.Second)
    con.AddSubscription([]string{"AMZN"})

	con.Wait()
}
