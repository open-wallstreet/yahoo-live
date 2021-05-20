# Yahoo Live Go Websocket

This package allow you to connect and retrieve websocket updates from Yahoo Finance tickers.

NOTE: Delay may wary and but can be a few seconds after real execution time. Should therefor not be used in applications that require fast and latest messages from the exchange.

## Installation

```bash
go get github.com/open-wallstreet/yahoo-live
```

### Example usage

See `examples` folder for more info

```go
import (
	"fmt"
	"time"

	"github.com/open-wallstreet/yahoo-live/pkg/yahoo"
	"github.com/open-wallstreet/yahoo-live/proto"
	"go.uber.org/zap"
)

func main() {
	logger, _ = zap.NewDevelopment()
    con, err := yahoo.NewWebsocket(logger.Sugar(), []string{"KIND-SDB.ST"})
	if err != nil {
		panic(err)
	}
	con.AddMessageHandler(on_msg)
	con.Wait()
}
func on_msg(message *proto.Yaticker) {
	println(fmt.Sprintf("%s: %s", time.Unix(message.Time/1000, 0).String(), message.String()))
}
```
