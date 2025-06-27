package cryptocurrency

import (
	"net/http"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Cryptocurrency struct {
	apiKey string
	doer   HTTPDoer
}

func New(apiKey string, doer HTTPDoer) *Cryptocurrency {
	return &Cryptocurrency{
		apiKey: apiKey,
		doer:   doer,
	}
}
