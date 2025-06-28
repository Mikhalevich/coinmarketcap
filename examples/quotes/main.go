package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Mikhalevich/coinmarketcap"
	"github.com/Mikhalevich/coinmarketcap/cryptocurrency"
	"github.com/Mikhalevich/coinmarketcap/currency"
)

const (
	timeout = time.Second * 5
	btcID   = 1
	ltcID   = 2
	usdID   = 2781
)

func main() {
	var (
		client = http.Client{
			Timeout: timeout,
		}

		prodExecutor = coinmarketcap.ProductionExecutor[*cryptocurrency.QuotesLatestResponse](
			os.Getenv("COIN_MARKET_CAP_KEY"),
			&client,
		)

		cc  = cryptocurrency.New(prodExecutor)
		log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	)

	rsp, err := cc.QuotesLatest(
		context.Background(),
		[]currency.Currency{currency.FromID(btcID), currency.FromID(ltcID)},
		[]currency.Currency{currency.FromID(usdID)},
	)

	if err != nil {
		log.Error("request quotes latest", "error", err.Error())
		os.Exit(1)
	}

	bytes, err := json.MarshalIndent(rsp, "", "	")
	if err != nil {
		log.Error("marshal json", "error", err.Error())
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, string(bytes))
}
