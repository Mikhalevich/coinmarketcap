package cryptocurrency

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Mikhalevich/coinmarketcap"
	"github.com/Mikhalevich/coinmarketcap/api/types"
	"github.com/Mikhalevich/coinmarketcap/currency"
)

const (
	infoEndpoint = "/v2/cryptocurrency/info"
)

type InfoResponse struct {
	Data   map[string]InfoData `json:"data"`
	Status types.Status        `json:"status"`
}

type InfoData struct {
	ID                            int              `json:"id"`
	Name                          string           `json:"name"`
	Symbol                        string           `json:"symbol"`
	Category                      string           `json:"category"`
	Slug                          string           `json:"slug"`
	Logo                          string           `json:"logo"`
	Description                   string           `json:"description"`
	DateAdded                     time.Time        `json:"date_added"`
	DateLaunched                  time.Time        `json:"date_launched"`
	Notice                        string           `json:"notice"`
	Tags                          []any            `json:"tags"`
	Platform                      types.PlatformV2 `json:"platform"`
	SelfReportedCirculatingSupply float64          `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         float64          `json:"self_reported_market_cap"`
	SelfReportedTags              []any            `json:"self_reported_tags"`
	InfiniteSupply                bool             `json:"infinite_supply"`
	Urls                          InfoUrls         `json:"urls"`
}

type InfoUrls struct {
	Website      []string `json:"website"`
	TechicalDoc  []string `json:"technical_doc"`
	Explored     []string `json:"explorer"`
	SourceCode   []string `json:"source_code"`
	MessageBoard []string `json:"message_board"`
	Chat         []string `json:"chat"`
	Announcement []string `json:"announcement"`
	Reddit       []string `json:"reddit"`
	Twitter      []string `json:"twitter"`
}

type infoOptions struct {
	Address     string
	Aux         []string
	SkipInvalid bool
}

// InfoOption info optional param.
type InfoOption func(opts *infoOptions)

// WithInfoAddress contract address.
func WithInfoAddress(address string) InfoOption {
	return func(opts *infoOptions) {
		opts.Address = address
	}
}

// WithInfoSkipInvalid specify request validation rules.
// When requesting records on multiple cryptocurrencies an error is returned
// if no match is found for 1 or more requested cryptocurrencies.
// If set to true, invalid lookups will be skipped allowing valid cryptocurrencies to still be returned.
// By default true.
func WithInfoSkipInvalid(skip bool) InfoOption {
	return func(opts *infoOptions) {
		opts.SkipInvalid = skip
	}
}

// WithInfoAux specify a list of supplemental data fields to return.
// By default "urls,logo,description,tags,platform,date_added,notice".
func WithInfoAux(fields ...string) InfoOption {
	return func(opts *infoOptions) {
		opts.Aux = fields
	}
}

// Info returns all static metadata available for one or more cryptocurrencies.
// https://coinmarketcap.com/api/documentation/v1/#operation/getV2CryptocurrencyInfo
func (c *Cryptocurrency) Info(
	ctx context.Context,
	currencies []currency.Currency,
	withOpts ...InfoOption,
) (*InfoResponse, error) {
	var (
		options = infoOptions{
			SkipInvalid: true,
		}
		response InfoResponse
	)

	for _, option := range withOpts {
		option(&options)
	}

	if err := c.executor.Get(
		ctx,
		infoEndpoint,
		func(req *http.Request) error {
			req.URL.RawQuery = makeInfoQuery(currencies, options)

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

func makeInfoQuery(
	currencies []currency.Currency,
	options infoOptions,
) string {
	query := make(url.Values)

	if options.Address != "" {
		query.Add("address", options.Address)
	} else {
		query.Add(makeCurrencyQueryKey(currencies), makeCommaSeparatedValues(convertCurrenciesToQueryKey(currencies)))
	}

	if len(options.Aux) > 0 {
		query.Add("aux", makeCommaSeparatedValues(options.Aux))
	}

	query.Add("skip_invalid", boolToString(options.SkipInvalid))

	return query.Encode()
}
