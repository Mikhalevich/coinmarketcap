package cryptocurrency

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Mikhalevich/coinmarketcap/currency"
)

const (
	baseURL = "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest"
	comma   = ","
)

type Money float64

type Quote struct {
	Price float64 `json:"price"`
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
	LastUpdated                   time.Time        `json:"last_updated"`
	SelfReportedCirculatingSupply float64          `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         float64          `json:"self_reported_market_cap"`
	Quotes                        map[string]Quote `json:"quote"`
}

type Status struct {
	Timestamp    time.Time `json:"timestamp"`
	ErrorCode    int       `json:"error_code"`
	ErrorMessage string    `json:"error_message"`
	Elapsed      int       `json:"elapsed"`
	CreditCount  int       `json:"credit_count"`
	Notice       string    `json:"notice"`
}

type QuotesLatestResponse struct {
	Data   map[string]Data `json:"data"`
	Status Status          `json:"status"`
}

func (q *QuotesLatestResponse) Quotes(toID currency.ID) map[currency.ID]Money {
	if len(q.Data) == 0 {
		return nil
	}

	quotes := make(map[currency.ID]Money, len(q.Data))

	for extID, data := range q.Data {
		quotes[currency.ID(extID)] = quotePrice(data, toID)
	}

	return quotes
}

func quotePrice(data Data, id currency.ID) Money {
	for extID, quote := range data.Quotes {
		if extID == id.String() {
			if quote.Price > 0 {
				return Money(1 / quote.Price)
			}

			return 0
		}
	}

	return 0
}

func (q *QuotesLatestResponse) IsError() bool {
	return q.Status.ErrorCode != 0
}

func (q *QuotesLatestResponse) ErrorMessage() string {
	return fmt.Sprintf("code: %d: msg: %s", q.Status.ErrorCode, q.Status.ErrorMessage)
}

// QuotesLatest returns the latest market quote for 1 or more cryptocurrencies.
// https://coinmarketcap.com/api/documentation/v1/#operation/getV2CryptocurrencyQuotesLatest
func (c *Cryptocurrency) QuotesLatest(
	ctx context.Context,
	convertFrom []currency.Currency,
	convertTo []currency.Currency,
) (*QuotesLatestResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	//nolint:canonicalheader
	req.Header.Set("X-CMC_PRO_API_KEY", c.apiKey)

	req.URL.RawQuery = makeQuery(convertFrom, convertTo)

	rsp, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do http request: %w", err)
	}

	defer rsp.Body.Close()

	var quotes QuotesLatestResponse
	if err := json.NewDecoder(rsp.Body).Decode(&quotes); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	if quotes.IsError() {
		return nil, fmt.Errorf("response error: %s", quotes.ErrorMessage())
	}

	return &quotes, nil
}

func makeQuery(from []currency.Currency, to []currency.Currency) string {
	q := make(url.Values)
	q.Add(makeConvertToQueryKey(to), makeCommaSeparatedValues(makeValues(to)))
	q.Add(makeConvertFromQueryKey(from), makeCommaSeparatedValues(makeValues(from)))

	return q.Encode()
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
