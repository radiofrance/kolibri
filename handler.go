package kolibri

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/xerrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	"github.com/radiofrance/kolibri/kind"
)

// UpdateHandlerPolicy that defines when two kubernetes objects are different.
type UpdateHandlerPolicy func(old, new metav1.Object) bool

// DefaultUpdateHandlerPolicy returns true only when the resource version changed.
func DefaultUpdateHandlerPolicy(old, new metav1.Object) bool {
	return old.GetResourceVersion() != new.GetResourceVersion()
}

// AlwaysUpdateHandlerPolicy always returns true.
func AlwaysUpdateHandlerPolicy(old, new metav1.Object) bool { return true }

// Handler is the main object of Kolibri. Like http.Handler, this object is used
// to "handle" events which occurs on the kubernetes cluster.
// TODO: Add more information
type Handler struct {
	ktx *Kontext

	events       eventRegistry
	updatePolicy UpdateHandlerPolicy

	kind     kind.Kind
	informer kind.Informer

	queue    workqueue.RateLimitingInterface
	recorder record.EventRecorder
}

func (ktr *Kontroller) NewHandler(knd kind.Kind, opts ...Option) (*Handler, error) {
	ctx := &handlerBuildContext{}

	if knd == nil {
		return nil, xerrors.New("knd can't be nil")
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
		ktx:      ktr.context(fmt.Sprintf("kolibri::%s::%s@%s", ktr.name, kind.FullName(knd), uuid.New().String())),
		kind:     knd,
		informer: knd.Informer(0, ctx.informerOpts...),
	}
	handler.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    handler.addHandler,
		UpdateFunc: handler.updateHandler,
		DeleteFunc: handler.deleteHandler,
	})

	if err := ctx.hdlrOpts.apply(handler); err != nil {
		return nil, err
	}
	if err := ctx.eventOpts.apply(&handler.events); err != nil {
		return nil, err
	}

	return handler, nil
}
