package cryptocurrency

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Mikhalevich/coinmarketcap"
	"github.com/Mikhalevich/coinmarketcap/api/types"
)

const (
	mapEndpoint = "/v1/cryptocurrency/map"
)

type MapResponse struct {
	Data   []MapData    `json:"data"`
	Status types.Status `json:"status"`
}

type MapStatus string

func (m MapStatus) String() string {
	return string(m)
}

const (
	MapStatusActive    MapStatus = "active"
	MapStatusInactive  MapStatus = "inactive"
	MapStatusUntracket MapStatus = "untracked"
)

type MapData struct {
	ID                  int            `json:"id"`
	Rank                float64        `json:"rank"`
	Name                string         `json:"name"`
	Symbol              string         `json:"symbol"`
	Slug                string         `json:"slug"`
	IsActive            int            `json:"is_active"`
	Status              MapStatus      `json:"status"`
	FirstHistoricalData time.Time      `json:"first_historical_data"`
	LastHistoricalData  time.Time      `json:"last_historical_data"`
	Platform            types.Platform `json:"platform"`
}

type MapSortField string

func (m MapSortField) String() string {
	return string(m)
}

const (
	MapSortID      MapSortField = "id"
	MapSortCMCRank MapSortField = "cmc_rank"
)

type mapOptions struct {
	ListingStatus MapStatus
	Start         int
	Limit         int
	Sort          MapSortField
	Symbol        []string
	Aux           []string
}

// MapOption map optional param.
type MapOption func(opts *mapOptions)

// WithListingStatus map satatus option.
// Default "active".
func WithMapListingStatus(status MapStatus) MapOption {
	return func(opts *mapOptions) {
		opts.ListingStatus = status
	}
}

// WithMapStart offset the start (1-based index) of the paginated list of items to return.
// Default 1.
func WithMapStart(start int) MapOption {
	return func(opts *mapOptions) {
		opts.Start = start
	}
}

// WithMapLimit use this parameter and the "start" parameter to determine your own pagination size.
func WithMapLimit(limit int) MapOption {
	return func(opts *mapOptions) {
		opts.Limit = limit
	}
}

// WithMapSort what field to sort the list of cryptocurrencies by.
// Default "id".
func WithMapSort(field MapSortField) MapOption {
	return func(opts *mapOptions) {
		opts.Sort = field
	}
}

// WithMapSymbol list of cryptocurrency symbols to return CoinMarketCap IDs for.
// If this option is passed, other options will be ignored.
func WithMapSymbol(symbols []string) MapOption {
	return func(opts *mapOptions) {
		opts.Symbol = symbols
	}
}

// WithMapAux specify a list of supplemental data fields to return.
// By default "platform,first_historical_data,last_historical_data,is_active".
func WithMapAux(fields ...string) MapOption {
	return func(opts *mapOptions) {
		opts.Aux = fields
	}
}

// Map returns a mapping of all cryptocurrencies to unique CoinMarketCap ids.
// https://coinmarketcap.com/api/documentation/v1/#operation/getV1CryptocurrencyMap
func (c *Cryptocurrency) Map(
	ctx context.Context,
	withOpts ...MapOption,
) (*MapResponse, error) {
	var (
		options = mapOptions{
			ListingStatus: "active",
			Start:         1,
			Sort:          MapSortID,
		}
		response MapResponse
	)

	for _, option := range withOpts {
		option(&options)
	}

	if err := c.executor.Get(
		ctx,
		mapEndpoint,
		func(req *http.Request) error {
			req.URL.RawQuery = makeMapQuery(options)

			return nil
		},
		&response,
	); err != nil {
		return nil, fmt.Errorf("execute get request: %w", err)
	}

	if response.Status.IsError() {
		return nil, coinmarketcap.NewError(response.Status.ErrorCode, response.Status.ErrorMessage)
	}

	return &response, nil
}

func makeMapQuery(options mapOptions) string {
	query := make(url.Values)

	if len(options.Symbol) > 0 {
		query.Add("symbol", makeCommaSeparatedValues(options.Symbol))

		return query.Encode()
	}

	if options.ListingStatus != "" {
		query.Add("listing_status", options.ListingStatus.String())
	}

	if options.Start > 0 {
		query.Add("start", strconv.Itoa(options.Start))
	}

	if options.Limit > 0 {
		query.Add("limit", strconv.Itoa(options.Limit))
	}

	if options.Sort != "" {
		query.Add("sort", options.Sort.String())
	}

	if len(options.Aux) > 0 {
		query.Add("aux", makeCommaSeparatedValues(options.Aux))
	}

	return query.Encode()
}
