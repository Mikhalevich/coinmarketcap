package key

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

type Key struct {
	executor Executor
}

func New(executor Executor) *Key {
	return &Key{
		executor: executor,
	}
}
