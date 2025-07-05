package coinmarketcap_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/Mikhalevich/coinmarketcap"
	"github.com/Mikhalevich/coinmarketcap/api/cryptocurrency"
	"github.com/stretchr/testify/require"
)

const (
	successResponseBody = `
{
	"data": {
			"1": {
					"id": 1,
					"name": "Bitcoin",
					"symbol": "BTC",
					"slug": "bitcoin",
					"is_active": 1,
					"is_fiat": 0,
					"cmc_rank": 1,
					"num_market_pairs": 12230,
					"circulating_supply": 19884678,
					"total_supply": 19884678,
					"market_cap_by_total_supply": 0,
					"max_supply": 21000000,
					"date_added": "2010-07-13T00:00:00Z",
					"tags": [
							{
									"category": "OTHERS",
									"name": "Mineable",
									"slug": "mineable"
							},
							{
									"category": "ALGORITHM",
									"name": "PoW",
									"slug": "pow"
							},
							{
									"category": "ALGORITHM",
									"name": "SHA-256",
									"slug": "sha-256"
							},
							{
									"category": "CATEGORY",
									"name": "Store Of Value",
									"slug": "store-of-value"
							},
							{
									"category": "CATEGORY",
									"name": "State Channel",
									"slug": "state-channel"
							},
							{
									"category": "CATEGORY",
									"name": "Coinbase Ventures Portfolio",
									"slug": "coinbase-ventures-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Three Arrows Capital Portfolio",
									"slug": "three-arrows-capital-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Polychain Capital Portfolio",
									"slug": "polychain-capital-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "YZi Labs Portfolio",
									"slug": "binance-labs-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Blockchain Capital Portfolio",
									"slug": "blockchain-capital-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "BoostVC Portfolio",
									"slug": "boostvc-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "CMS Holdings Portfolio",
									"slug": "cms-holdings-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "DCG Portfolio",
									"slug": "dcg-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "DragonFly Capital Portfolio",
									"slug": "dragonfly-capital-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Electric Capital Portfolio",
									"slug": "electric-capital-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Fabric Ventures Portfolio",
									"slug": "fabric-ventures-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Framework Ventures Portfolio",
									"slug": "framework-ventures-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Galaxy Digital Portfolio",
									"slug": "galaxy-digital-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Huobi Capital Portfolio",
									"slug": "huobi-capital-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Alameda Research Portfolio",
									"slug": "alameda-research-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "a16z Portfolio",
									"slug": "a16z-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "1Confirmation Portfolio",
									"slug": "1confirmation-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Winklevoss Capital Portfolio",
									"slug": "winklevoss-capital-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "USV Portfolio",
									"slug": "usv-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Placeholder Ventures Portfolio",
									"slug": "placeholder-ventures-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Pantera Capital Portfolio",
									"slug": "pantera-capital-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Multicoin Capital Portfolio",
									"slug": "multicoin-capital-portfolio"
							},
							{
									"category": "CATEGORY",
									"name": "Paradigm Portfolio",
									"slug": "paradigm-portfolio"
							},
							{
									"category": "PLATFORM",
									"name": "Bitcoin Ecosystem",
									"slug": "bitcoin-ecosystem"
							},
							{
									"category": "CATEGORY",
									"name": "Layer 1",
									"slug": "layer-1"
							},
							{
									"category": "CATEGORY",
									"name": "FTX Bankruptcy Estate ",
									"slug": "ftx-bankruptcy-estate"
							},
							{
									"category": "CATEGORY",
									"name": "2017/18 Alt season",
									"slug": "2017-2018-alt-season"
							},
							{
									"category": "CATEGORY",
									"name": "US Strategic Crypto Reserve",
									"slug": "us-strategic-crypto-reserve"
							},
							{
									"category": "CATEGORY",
									"name": "Binance Ecosystem",
									"slug": "binance-ecosystem"
							},
							{
									"category": "CATEGORY",
									"name": "Binance Listing",
									"slug": "binance-listing"
							}
					],
					"platform": {
							"id": "",
							"name": "",
							"symbol": "",
							"slug": "",
							"token_address": ""
					},
					"last_updated": "2025-06-28T16:19:00Z",
					"self_reported_circulating_supply": 0,
					"self_reported_market_cap": 0,
					"quote": {
							"2781": {
									"price": 107431.72910155341,
									"volume_24h": 33308615293.517628,
									"volume_change_24h": -24.3035,
									"volume_24h_reported": 0,
									"volume_7d": 0,
									"volume_7d_reported": 0,
									"volume_30d": 0,
									"market_cap": 2136245340167.619,
									"market_cap_dominance": 64.7794,
									"fully_diluted_market_cap": 2256066311132.62,
									"percent_change_1h": 0.13970253,
									"percent_change_24h": 0.16709607,
									"percent_change_7d": 3.67255782,
									"percent_change_30d": 0.37710465,
									"last_updated": "2025-06-28T16:19:00Z"
							}
					}
			},
			"2": {
					"id": 2,
					"name": "Litecoin",
					"symbol": "LTC",
					"slug": "litecoin",
					"is_active": 1,
					"is_fiat": 0,
					"cmc_rank": 20,
					"num_market_pairs": 1377,
					"circulating_supply": 76016051.98347135,
					"total_supply": 84000000,
					"market_cap_by_total_supply": 0,
					"max_supply": 84000000,
					"date_added": "2013-04-28T00:00:00Z",
					"tags": [
							{
									"category": "OTHERS",
									"name": "Mineable",
									"slug": "mineable"
							},
							{
									"category": "ALGORITHM",
									"name": "PoW",
									"slug": "pow"
							},
							{
									"category": "ALGORITHM",
									"name": "Scrypt",
									"slug": "scrypt"
							},
							{
									"category": "INDUSTRY",
									"name": "Medium of Exchange",
									"slug": "medium-of-exchange"
							},
							{
									"category": "CATEGORY",
									"name": "Layer 1",
									"slug": "layer-1"
							},
							{
									"category": "CATEGORY",
									"name": "2017/18 Alt season",
									"slug": "2017-2018-alt-season"
							},
							{
									"category": "CATEGORY",
									"name": "Made in America",
									"slug": "made-in-america"
							},
							{
									"category": "CATEGORY",
									"name": "Binance Ecosystem",
									"slug": "binance-ecosystem"
							},
							{
									"category": "CATEGORY",
									"name": "Binance Listing",
									"slug": "binance-listing"
							}
					],
					"platform": {
							"id": "",
							"name": "",
							"symbol": "",
							"slug": "",
							"token_address": ""
					},
					"last_updated": "2025-06-28T16:19:00Z",
					"self_reported_circulating_supply": 0,
					"self_reported_market_cap": 0,
					"quote": {
							"2781": {
									"price": 86.61639053394676,
									"volume_24h": 272507551.2196314,
									"volume_change_24h": -16.0333,
									"volume_24h_reported": 0,
									"volume_7d": 0,
									"volume_7d_reported": 0,
									"volume_30d": 0,
									"market_cap": 6584236045.449152,
									"market_cap_dominance": 0.1997,
									"fully_diluted_market_cap": 7275776804.85,
									"percent_change_1h": 0.88738713,
									"percent_change_24h": 2.85256281,
									"percent_change_7d": 4.36907952,
									"percent_change_30d": -9.76764849,
									"last_updated": "2025-06-28T16:19:00Z"
							}
					}
			}
	},
	"status": {
			"timestamp": "2025-06-28T16:19:48.947Z",
			"error_code": 0,
			"error_message": "",
			"elapsed": 19,
			"credit_count": 1,
			"notice": ""
	}
}
	`
)

func BenchmarkRequestExecutor(b *testing.B) {
	var (
		ctrl     = gomock.NewController(b)
		doer     = coinmarketcap.NewMockHTTPDoer(ctrl)
		executor = coinmarketcap.NewRequestExecutor("testApiKey", "testHost", doer)
	)

	for b.Loop() {
		doer.EXPECT().
			Do(gomock.Any()).
			Return(
				&http.Response{
					Body: io.NopCloser(strings.NewReader(successResponseBody)),
				},
				nil,
			)

		var rsp cryptocurrency.QuotesLatestResponse
		err := executor.Get(
			b.Context(),
			"some_path",
			func(req *http.Request) error {
				return nil
			},
			&rsp,
		)

		if err != nil {
			b.Fatalf("unexpected get error: %v", err)
		}

		if len(rsp.Data) == 0 {
			b.Fatal("invalid response data")
		}
	}
}

func TestMakeEndpointError(t *testing.T) {
	t.Parallel()

	var (
		ctrl = gomock.NewController(t)
		doer = coinmarketcap.NewMockHTTPDoer(ctrl)
	)

	t.Run("missing protocol", func(t *testing.T) {
		t.Parallel()

		executor := coinmarketcap.NewRequestExecutor("testApiKey", ":some_host", doer)

		var rsp cryptocurrency.InfoResponse
		err := executor.Get(
			t.Context(),
			"some_path",
			func(req *http.Request) error {
				return nil
			},
			&rsp,
		)

		require.EqualError(t, err, "make endpoint url: parse \":some_host\": missing protocol scheme")
	})
}

func TestMakeRequestError(t *testing.T) {
	t.Parallel()

	var (
		ctrl     = gomock.NewController(t)
		doer     = coinmarketcap.NewMockHTTPDoer(ctrl)
		executor = coinmarketcap.NewRequestExecutor("testApiKey", "some_host", doer)
	)

	var rsp cryptocurrency.InfoResponse
	//nolint:staticcheck
	err := executor.Get(
		nil,
		"some_path",
		func(req *http.Request) error {
			return errors.New("some pre process error")
		},
		&rsp,
	)

	require.EqualError(t, err, "create http request: net/http: nil Context")
}

func TestPreProccessError(t *testing.T) {
	t.Parallel()

	var (
		ctrl     = gomock.NewController(t)
		doer     = coinmarketcap.NewMockHTTPDoer(ctrl)
		executor = coinmarketcap.NewRequestExecutor("testApiKey", "some_host", doer)
	)

	var rsp cryptocurrency.InfoResponse
	err := executor.Get(
		t.Context(),
		"some_path",
		func(req *http.Request) error {
			return errors.New("some pre process error")
		},
		&rsp,
	)

	require.EqualError(t, err, "pre process: some pre process error")
}

func TestDoerError(t *testing.T) {
	t.Parallel()

	var (
		ctrl     = gomock.NewController(t)
		doer     = coinmarketcap.NewMockHTTPDoer(ctrl)
		executor = coinmarketcap.NewRequestExecutor("testApiKey", "some_host", doer)
	)

	doer.EXPECT().Do(gomock.Any()).Return(nil, errors.New("some do error"))

	var rsp cryptocurrency.InfoResponse
	err := executor.Get(
		t.Context(),
		"some_path",
		func(req *http.Request) error {
			return nil
		},
		&rsp,
	)

	require.EqualError(t, err, "do http request: some do error")
}

func TestDecodeResponseJSONError(t *testing.T) {
	t.Parallel()

	var (
		ctrl     = gomock.NewController(t)
		doer     = coinmarketcap.NewMockHTTPDoer(ctrl)
		executor = coinmarketcap.NewRequestExecutor("testApiKey", "some_host", doer)
	)

	t.Run("empty response", func(t *testing.T) {
		t.Parallel()

		doer.EXPECT().
			Do(gomock.Any()).
			Return(&http.Response{
				Body: io.NopCloser(strings.NewReader("")),
			}, nil)

		var rsp cryptocurrency.InfoResponse
		err := executor.Get(
			t.Context(),
			"some_path",
			func(req *http.Request) error {
				return nil
			},
			&rsp,
		)

		require.EqualError(t, err, "json decode: EOF")
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		doer.EXPECT().
			Do(gomock.Any()).
			Return(&http.Response{
				Body: io.NopCloser(strings.NewReader("invalid json")),
			}, nil)

		var rsp cryptocurrency.InfoResponse
		err := executor.Get(
			t.Context(),
			"some_path",
			func(req *http.Request) error {
				return nil
			},
			&rsp,
		)

		require.EqualError(t, err, "json decode: invalid character 'i' looking for beginning of value")
	})

	t.Run("json array instead of object", func(t *testing.T) {
		t.Parallel()

		doer.EXPECT().
			Do(gomock.Any()).
			Return(&http.Response{
				Body: io.NopCloser(strings.NewReader("[]")),
			}, nil)

		var rsp cryptocurrency.InfoResponse
		err := executor.Get(
			t.Context(),
			"some_path",
			func(req *http.Request) error {
				return nil
			},
			&rsp,
		)

		require.EqualError(t, err,
			"json decode: json: cannot unmarshal array into Go value of type cryptocurrency.InfoResponse")
	})
}

func TestSuccess(t *testing.T) {
	t.Parallel()

	var (
		ctrl     = gomock.NewController(t)
		doer     = coinmarketcap.NewMockHTTPDoer(ctrl)
		executor = coinmarketcap.NewRequestExecutor("testApiKey", "some_host", doer)
	)

	doer.EXPECT().
		Do(gomock.Any()).
		Return(
			&http.Response{
				Body: io.NopCloser(strings.NewReader(successResponseBody)),
			},
			nil,
		)

	var rsp cryptocurrency.InfoResponse
	err := executor.Get(
		t.Context(),
		"some_path",
		func(req *http.Request) error {
			return nil
		},
		&rsp,
	)

	require.NoError(t, err)
}
