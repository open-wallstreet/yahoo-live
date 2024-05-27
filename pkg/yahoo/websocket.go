package yahoo

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/fasthttp/websocket"
	yaproto "github.com/open-wallstreet/yahoo-live/pkg/yahoo/proto"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const url = "wss://streamer.finance.yahoo.com/"

type Subscription struct {
	Subscribe []string `json:"subscribe"`
}

type SocketConnection struct {
	logger              *zap.SugaredLogger
	interruptSignal     chan os.Signal
	done                chan struct{}
	on_message_handlers []func(message *yaproto.Yaticker)
	connection          *websocket.Conn
}

// Add new message handler function that will be called for all new messages websocket receives
func (s *SocketConnection) AddMessageHandler(f func(message *yaproto.Yaticker)) {
	s.on_message_handlers = append(s.on_message_handlers, f)
}

// Wait runs infinite for loop until a interuptSignal or close singal
// has been received by the message handlers. Will excited after running Close()
func (s *SocketConnection) Wait() {
	for {
		select {
		case <-s.done:
			return
		case <-s.interruptSignal:
			log.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := s.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				s.logger.Errorf("write close:", err)
				return
			}
			<-s.done
			return
		}

	}
}

// Closes websocket connection and send done signal for gracefull shutdown
func (s *SocketConnection) Close() {
	s.connection.Close()
	close(s.done)
}

func (s *SocketConnection) handleMessages() {
	defer s.Close()
	for {
		_, message, err := s.connection.ReadMessage()
		if err != nil {
			switch err.(type) {
			case *websocket.CloseError:
				s.logger.Info("received close message")
				return
			default:
				s.logger.Errorf("read: %v", err)
				return
			}

		}
		msg, err := Base64Decode(message)
		if err != nil {
			s.logger.Errorf("failed to decode: %v", err)
			return
		}
		ticker := &yaproto.Yaticker{}
		if err := proto.Unmarshal(msg, ticker); err != nil {
			s.logger.Errorf("failed to parse %v", msg)
		}
		for _, handler := range s.on_message_handlers {
			handler(ticker)
		}
	}
}

// Creates a new yahoo.SocketConnection and subscribe to all tickers listed in array
func NewWebsocket(logger *zap.SugaredLogger, tickers []string) (*SocketConnection, error) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Panicf("failed to connect to websocket %v", err)
	}
	socketConnection := &SocketConnection{
		logger:              logger,
		interruptSignal:     interrupt,
		on_message_handlers: []func(message *yaproto.Yaticker){},
		connection:          connection,
		done:                make(chan struct{}),
	}

	go socketConnection.handleMessages()

	msg := Subscription{
		Subscribe: tickers,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		logger.Errorf("failed to marshal subscription message:", err)
		return nil, err
	}
	err = connection.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		logger.Infof("failed to send subscription message:", err)
		return nil, err
	}
	return socketConnection, nil
}


func (s *SocketConnection) AddSubscription(tickers []string) {
    msg := Subscription{
        Subscribe: tickers,
    }
    b, err := json.Marshal(msg)
    if err != nil {
        s.logger.Errorf("Failed to Marshal Subscription Message", err)
    }

    err = s.connection.WriteMessage(websocket.TextMessage, b)

    if err != nil  {
        s.logger.Errorf("Failed to send subscription ", err)
    }
}

func Base64Decode(str []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(str))
}
