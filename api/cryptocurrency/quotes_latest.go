package cryptocurrency

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Mikhalevich/coinmarketcap"
	"github.com/Mikhalevich/coinmarketcap/currency"
)

const (
	quoteLatestEndpoint = "/v2/cryptocurrency/quotes/latest"
	comma               = ","
)

type QuotesLatestResponse struct {
	Data   map[string]Data `json:"data"`
	Status Status          `json:"status"`
}

type Data struct {
	ID                            int              `json:"id"`
	Name                          string           `json:"name"`
	Symbol                        string           `json:"symbol"`
	Slug                          string           `json:"slug"`
	IsActive                      int              `json:"is_active"`
	IsFiat                        int              `json:"is_fiat"`
	CMCRank                       int              `json:"cmc_rank"`
	NumMarketPairs                int              `json:"num_market_pairs"`
	CirculatingSupply             float64          `json:"circulating_supply"`
	TotalSupply                   float64          `json:"total_supply"`
	MarketCapByTotalSupply        float64          `json:"market_cap_by_total_supply"`
	MaxSupply                     float64          `json:"max_supply"`
	DateAdded                     time.Time        `json:"date_added"`
	Tags                          []any            `json:"tags"`
	Platform                      Platform         `json:"platform"`
	LastUpdated                   time.Time        `json:"last_updated"`
	SelfReportedCirculatingSupply float64          `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         float64          `json:"self_reported_market_cap"`
	Quotes                        map[string]Quote `json:"quote"`
}

type Platform struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Slug        string `json:"slug"`
	TokenAdress string `json:"token_address"`
}

type Quote struct {
	Price                 float64   `json:"price"`
	Volume24h             float64   `json:"volume_24h"`
	VolumeChange24h       float64   `json:"volume_change_24h"`
	Volume24hReported     float64   `json:"volume_24h_reported"`
	Volume7d              float64   `json:"volume_7d"`
	Volume7dReported      float64   `json:"volume_7d_reported"`
	Volume30d             float64   `json:"volume_30d"`
	MarketCap             float64   `json:"market_cap"`
	MarketCapDominance    float64   `json:"market_cap_dominance"`
	FullyDilutedMarketCap float64   `json:"fully_diluted_market_cap"`
	PercentChange1h       float64   `json:"percent_change_1h"`
	PercentChange24h      float64   `json:"percent_change_24h"`
	PercentChange7d       float64   `json:"percent_change_7d"`
	PercentChange30d      float64   `json:"percent_change_30d"`
	LastUpdated           time.Time `json:"last_updated"`
}

type Status struct {
	Timestamp    time.Time `json:"timestamp"`
	ErrorCode    int       `json:"error_code"`
	ErrorMessage string    `json:"error_message"`
	Elapsed      int       `json:"elapsed"`
	CreditCount  int       `json:"credit_count"`
	Notice       string    `json:"notice"`
}

func (q *QuotesLatestResponse) QuotePrices(baseSymbol string) map[string]float64 {
	if len(q.Data) == 0 {
		return nil
	}

	quotes := make(map[string]float64, len(q.Data))

	for symbol, data := range q.Data {
		quotes[symbol] = quotePrice(data, baseSymbol)
	}

	return quotes
}

func quotePrice(data Data, baseSymbol string) float64 {
	for symbol, quote := range data.Quotes {
		if symbol == baseSymbol {
			if quote.Price > 0 {
				return 1 / quote.Price
			}

			return 0
		}
	}

	return 0
}

func (q *QuotesLatestResponse) IsError() bool {
	return q.Status.ErrorCode != 0
}

type quotesLatestOptions struct {
	Aux         []string
	SkipInvalid bool
}

// QuotesLatestOption quotes latest optional param.
type QuotesLatestOption func(opts *quotesLatestOptions)

// WithSkipInvalid specify request validation rules.
// When requesting records on multiple cryptocurrencies an error is returned
// if no match is found for 1 or more requested cryptocurrencies.
// If set to true, invalid lookups will be skipped allowing valid cryptocurrencies to still be returned.
// By default true.
func WithSkipInvalid(skip bool) QuotesLatestOption {
	return func(opts *quotesLatestOptions) {
		opts.SkipInvalid = skip
	}
}

// WithAux specify a list of supplemental data fields to return.
// By default "num_market_pairs, cmc_rank, date_added, tags, platform,
// max_supply, circulating_supply, total_supply, is_active, is_fiat".
func WithAux(fields ...string) QuotesLatestOption {
	return func(opts *quotesLatestOptions) {
		opts.Aux = fields
	}
}

// QuotesLatest returns the latest market quote for 1 or more cryptocurrencies.
// https://coinmarketcap.com/api/documentation/v1/#operation/getV2CryptocurrencyQuotesLatest
func (c *Cryptocurrency) QuotesLatest(
	ctx context.Context,
	convertFrom []currency.Currency,
	convertTo []currency.Currency,
	withOpts ...QuotesLatestOption,
) (*QuotesLatestResponse, error) {
	var (
		options = quotesLatestOptions{
			SkipInvalid: true,
		}
		quotes QuotesLatestResponse
	)

	for _, option := range withOpts {
		option(&options)
	}

	if err := c.executor.Get(
		ctx,
		quoteLatestEndpoint,
		func(req *http.Request) error {
			req.URL.RawQuery = makeQuery(convertFrom, convertTo, options)

			return nil
		},
		&quotes,
	); err != nil {
		return nil, fmt.Errorf("execute get request: %w", err)
	}

	if quotes.IsError() {
		return nil, coinmarketcap.NewError(quotes.Status.ErrorCode, quotes.Status.ErrorMessage)
	}

	return &quotes, nil
}

func makeQuery(
	from []currency.Currency,
	to []currency.Currency,
	options quotesLatestOptions,
) string {
	query := make(url.Values)
	query.Add(makeConvertToQueryKey(to), makeCommaSeparatedValues(makeValues(to)))
	query.Add(makeConvertFromQueryKey(from), makeCommaSeparatedValues(makeValues(from)))

	if len(options.Aux) > 0 {
		query.Add("aux", makeCommaSeparatedValues(options.Aux))
	}

	query.Add("skip_invalid", boolToString(options.SkipInvalid))

	return query.Encode()
}

func boolToString(b bool) string {
	if b {
		return "true"
	}

	return "false"
}

func makeConvertFromQueryKey(from []currency.Currency) string {
	if len(from) == 0 {
		return ""
	}

	if from[0].ID != "" {
		return "id"
	}

	if from[0].Symbol != "" {
		return "symbol"
	}

	if from[0].Slug != "" {
		return "slug"
	}

	return ""
}

func makeConvertToQueryKey(from []currency.Currency) string {
	if len(from) == 0 {
		return ""
	}

	if from[0].ID != "" {
		return "convert_id"
	}

	if from[0].Symbol != "" {
		return "convert"
	}

	return ""
}

func makeValues(from []currency.Currency) []string {
	if len(from) == 0 {
		return nil
	}

	values := make([]string, 0, len(from))

	for _, curr := range from {
		if curr.ID != "" {
			values = append(values, curr.ID.String())

			continue
		}

		if curr.Symbol != "" {
			values = append(values, curr.Symbol)

			continue
		}

		if curr.Slug != "" {
			values = append(values, curr.Slug)
		}
	}

	return values
}

func makeCommaSeparatedValues(currencies []string) string {
	switch len(currencies) {
	case 0:
		return ""
	case 1:
		return currencies[0]
	}

	var bufLen int
	bufLen += len(comma)*len(currencies) - 1

	for _, curr := range currencies {
		bufLen += len(curr)
	}

	var builder strings.Builder

	builder.Grow(bufLen)

	builder.WriteString(currencies[0])

	for _, curr := range currencies[1:] {
		builder.WriteString(comma)
		builder.WriteString(curr)
	}

	return builder.String()
}
