package kolibri

import (
	"context"

	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
	"k8s.io/client-go/kubernetes"

	"github.com/radiofrance/kolibri/kind"
	"github.com/radiofrance/kolibri/log"
	"github.com/radiofrance/kolibri/log/fake"
)

type Kontroller struct {
	kubernetes.Interface

	name   string
	ctx    context.Context
	logger log.Logger

	handlers []*Handler
}

func NewController(name string, client kubernetes.Interface, opts ...interface{}) (*Kontroller, error) {
	if client == nil {
		return nil, xerrors.New("client cannot be nil")
	}

	return &Kontroller{
		Interface: client,
		name:      name,
		logger:    fake.New(),
	}, nil
}

func (ktr *Kontroller) SetLogger(logger log.Logger) {
	if logger != nil {
		ktr.logger = logger
	}
}
func (ktr *Kontroller) Register(handlers ...*Handler) error {
	ktr.handlers = append(ktr.handlers, handlers...)
	return nil
}
func (ktr *Kontroller) Run(ctx context.Context) error {
	errg, ctx := errgroup.WithContext(ctx)

	for _, handler := range ktr.handlers {
		errg.Go(func() error { return handler.Run(ctx, ktr) })
	}

	return errg.Wait()
}

func (ktr *Kontroller) context(name string, informer kind.Informer) *Kontext {
	return &Kontext{
		Context:  context.WithValue(ktr.ctx, KontextKey("name"), name),
		Logger:   ktr.logger.Named(name),
		informer: informer,
	}
}
