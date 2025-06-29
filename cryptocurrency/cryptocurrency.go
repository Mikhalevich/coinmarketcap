package cryptocurrency

import (
	"context"
	"net/http"
)

type Executor interface {
	Get(
		ctx context.Context,
		path string,
		preProcessFn func(req *http.Request) error,
		result any,
	) error
}

type Cryptocurrency struct {
	executor Executor
}

func New(executor Executor) *Cryptocurrency {
	return &Cryptocurrency{
		executor: executor,
	}
}
