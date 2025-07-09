package fiat

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Mikhalevich/coinmarketcap"
	"github.com/Mikhalevich/coinmarketcap/api/types"
)

const (
	mapEndpoint = "/v1/fiat/map"
)

type MapResponse struct {
	Data   []MapData    `json:"data"`
	Status types.Status `json:"status"`
}

type MapData struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Sign   string `json:"sign"`
	Symbol string `json:"symbol"`
}

type mapOptions struct {
	Start         int
	Limit         int
	Sort          string
	IncludeMetals bool
}

// MapOption optional map parameter.
type MapOption func(opts *mapOptions)

// WithMapStart offset the start (1-based index) of the paginated list of items to return.
// Default 1.
func WithMapStart(start int) MapOption {
	return func(opts *mapOptions) {
		opts.Start = start
	}
}

// WithMapLimit specify the number of results to return.
// Use this parameter and the "start" parameter to determine your own pagination size.
func WithMapLimit(limit int) MapOption {
	return func(opts *mapOptions) {
		opts.Limit = limit
	}
}

// WithMapSort what field to sort the list by.
// Default "id".
func WithMapSort(field string) MapOption {
	return func(opts *mapOptions) {
		opts.Sort = field
	}
}

// WithMapMetals pass true to include precious metals.
// Default false.
func WithMapMetals(include bool) MapOption {
	return func(opts *mapOptions) {
		opts.IncludeMetals = include
	}
}

// Map returns a mapping of all supported fiat currencies to unique CoinMarketCap ids.
// https://coinmarketcap.com/api/documentation/v1/#operation/getV1FiatMap
func (f *Fiat) Map(
	ctx context.Context,
	withOpts ...MapOption,
) (*MapResponse, error) {
	var (
		options = mapOptions{
			Start: 1,
		}
		response MapResponse
	)

	for _, option := range withOpts {
		option(&options)
	}

	if err := f.executor.Get(
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

func makeMapQuery(opts mapOptions) string {
	query := make(url.Values)

	if opts.Start > 1 {
		query.Add("start", strconv.Itoa(opts.Start))
	}

	if opts.Limit > 0 {
		query.Add("limit", strconv.Itoa(opts.Limit))
	}

	if opts.Sort != "" {
		query.Add("sort", opts.Sort)
	}

	if opts.IncludeMetals {
		query.Add("include_metals", strconv.FormatBool(opts.IncludeMetals))
	}

	return query.Encode()
}
