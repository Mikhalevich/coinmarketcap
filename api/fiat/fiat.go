package fiat

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

type Fiat struct {
	executor Executor
}

func New(executor Executor) *Fiat {
	return &Fiat{
		executor: executor,
	}
}
