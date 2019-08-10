package kolibri

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	"github.com/radiofrance/kolibri/kind"
)

// Handler is the main object of Kolibri. Like http.Handler, this object is used
// to "handle" events which occurs on the kubernetes cluster.
// TODO: Add more information
type Handler struct {
	ktr    *Kontroller
	events eventRegistry

	kind     kind.Kind
	informer kind.Informer

	queue    workqueue.RateLimitingInterface
	recorder record.EventRecorder
}

func (k *Kontroller) NewHandler(kind kind.Kind, opts ...Option) (*Handler, error) {
	ctx := &handlerBuildContext{}

	if kind == nil {
		return nil, xerrors.New("kind can't be nil")
	}

	for _, opt := range opts {
		if err := opt.apply(ctx); err != nil {
			return nil, err
		}
	}

	if len(ctx.eventOpts) == 0 {
		return nil, xerrors.Errorf("at least one event handler (On...) must be provided")
	}

	handler := &Handler{
		ktr:      k.copy(),
		kind:     kind,
		informer: kind.Informer(5*time.Second, ctx.informerOpts...),
		queue: workqueue.NewNamedRateLimitingQueue(
			workqueue.DefaultControllerRateLimiter(),
			fmt.Sprintf("%s:%s:%s/%s@%d", "kolibris", k.name, kind.APIVersion(), kind.Name(), uuid.New().String()),
		),
	}
	handler.ktr.Logger = k.Named(fmt.Sprintf("%s/%s", kind.APIVersion(), kind.Name()))
	handler.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    handler.addHandler,
		UpdateFunc: handler.updateHandler,
		DeleteFunc: handler.deleteHandler,
	})

	if err := ctx.ktrlOpts.apply(handler.ktr); err != nil {
		return nil, err
	}
	if err := ctx.eventOpts.apply(&handler.events); err != nil {
		return nil, err
	}

	return handler, nil
}
