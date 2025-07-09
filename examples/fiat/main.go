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
	"github.com/Mikhalevich/coinmarketcap/api/fiat"
)

const (
	timeout = time.Second * 5
)

func main() {
	var (
		client = http.Client{
			Timeout: timeout,
		}

		prodExecutor = coinmarketcap.ProductionExecutor(os.Getenv("COIN_MARKET_CAP_KEY"), &client)
		f            = fiat.New(prodExecutor)
		log          = slog.New(slog.NewTextHandler(os.Stdout, nil))
	)

	quotes, err := f.Map(
		context.Background(),
		fiat.WithMapSort("name"),
		fiat.WithMapMetals(true),
	)
	if err != nil {
		log.Error("request fiat map", "error", err.Error())
		os.Exit(1)
	}

	jsonPrint(quotes, log)
}

func jsonPrint(response any, log *slog.Logger) {
	bytes, err := json.MarshalIndent(response, "", "	")
	if err != nil {
		log.Error("marshal json", "error", err.Error())
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, string(bytes))
}
