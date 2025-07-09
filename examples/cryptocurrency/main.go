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
	"github.com/Mikhalevich/coinmarketcap/api/cryptocurrency"
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

		prodExecutor = coinmarketcap.ProductionExecutor(os.Getenv("COIN_MARKET_CAP_KEY"), &client)
		cryptoc      = cryptocurrency.New(prodExecutor)
		log          = slog.New(slog.NewTextHandler(os.Stdout, nil))
	)

	quotes, err := cryptoc.QuotesLatest(
		context.Background(),
		[]currency.Currency{currency.FromID(btcID), currency.FromID(ltcID)},
		[]currency.Currency{currency.FromID(usdID)},
		cryptocurrency.WithQLSkipInvalid(false),
	)
	if err != nil {
		log.Error("request quotes latest", "error", err.Error())
		os.Exit(1)
	}

	jsonPrint(quotes, log)

	info, err := cryptoc.Info(
		context.Background(),
		[]currency.Currency{currency.Slug("bitcoin"), currency.Slug("litecoin")},
	)
	if err != nil {
		log.Error("request info", "error", err.Error())
		os.Exit(1)
	}

	jsonPrint(info, log)

	mappings, err := cryptoc.Map(
		context.Background(),
		cryptocurrency.WithMapSymbol("BTC"),
	)
	if err != nil {
		log.Error("request map", "error", err.Error())
		os.Exit(1)
	}

	jsonPrint(mappings, log)
}

func jsonPrint(response any, log *slog.Logger) {
	bytes, err := json.MarshalIndent(response, "", "	")
	if err != nil {
		log.Error("marshal json", "error", err.Error())
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, string(bytes))
}
