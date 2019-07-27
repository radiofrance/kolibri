package kolibri

import (
	"context"
)

type Handler interface {
	Run(ctx context.Context) error
}
type Option interface{}

func NewHandler(opts ...Option) (Handler, error) { return nil, nil }
