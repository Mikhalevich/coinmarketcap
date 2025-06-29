package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	productionHost = "https://pro-api.coinmarketcap.com"
)

// ProductionExecutor constructs production request executor with https://pro-api.coinmarketcap.com base url.
func ProductionExecutor(apiKey string, doer HTTPDoer) *RequestExecutor {
	return NewRequestExecutor(apiKey, productionHost, doer)
}

// HTTPDoer interface for external implementation for doint http request.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// RequestExecutor structure for raw request executing for coinmarketcap api.
type RequestExecutor struct {
	apiKey string
	host   string
	doer   HTTPDoer
}

// NewRequestExecutor construct new request executor.
func NewRequestExecutor(apiKey string, host string, doer HTTPDoer) *RequestExecutor {
	return &RequestExecutor{
		apiKey: apiKey,
		host:   host,
		doer:   doer,
	}
}

// Get execute Get request for specified endpoint path.
// before executing request preProcessFn function is invoked with request object.
func (re *RequestExecutor) Get(
	ctx context.Context,
	path string,
	preProcessFn func(req *http.Request) error,
	result any,
) error {
	endpointURL, err := url.JoinPath(re.host, path)
	if err != nil {
		return fmt.Errorf("make endpoint url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return fmt.Errorf("create http request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	//nolint:canonicalheader
	req.Header.Set("X-CMC_PRO_API_KEY", re.apiKey)

	if err := preProcessFn(req); err != nil {
		return fmt.Errorf("pre process: %w", err)
	}

	rsp, err := re.doer.Do(req)
	if err != nil {
		return fmt.Errorf("do http request: %w", err)
	}

	defer rsp.Body.Close()

	if err := json.NewDecoder(rsp.Body).Decode(result); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	return nil
}
