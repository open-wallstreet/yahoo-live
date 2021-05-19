package yahoo

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	yaproto "github.com/open-wallstreet/yahoo-live/proto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const url = "wss://streamer.finance.yahoo.com/"

type Subscription struct {
	Subscribe []string `json:"subscribe"`
}

func Base64Decode(str []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(str))
}

func NewYahooWebsocket(logger *zap.SugaredLogger, tickers []string, onMessage func(message *yaproto.Yaticker)) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Panicf("failed to connect to websocket %v", err)
	}
	defer connection.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				logger.Errorf("read: %v", err)
				return
			}
			msg, err := Base64Decode(message)
			if err != nil {
				logger.Errorf("failed to decode: %v", err)
				return
			}
			ticker := &yaproto.Yaticker{}
			if err := proto.Unmarshal(msg, ticker); err != nil {
				logger.Errorf("failed to parse %v", msg)
			}
			// logger.Infof("recv: %v", ticker)
			onMessage(ticker)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	msg := Subscription{
		Subscribe: tickers,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("failed to marshal subscription message:", err)
		return
	}
	err = connection.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		logger.Infof("failed to send subscription message:", err)
		return
	}

	for {
		select {
		case <-done:
			return

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}

	}
}
