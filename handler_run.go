package kolibri

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/xerrors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

func (h *Handler) Run(ctx context.Context, ktr *Kontroller) error {
	h.queue = workqueue.NewNamedRateLimitingQueue(
		workqueue.DefaultControllerRateLimiter(),
		h.ktx.Value(KontextKey("name")).(string),
	)
	defer h.queue.ShutDown()

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(h.ktx.Debugf)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: ktr.CoreV1().Events("")})
	h.recorder = eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: h.ktx.Value(KontextKey("name")).(string)})

	h.informer.Start(ctx.Done())
	if ok := cache.WaitForCacheSync(ctx.Done(), h.informer.HasSynced); !ok {
		panic("failed to wait for caches to sync")
	}

	for i := 0; i < 10; i++ {
		go wait.Until(func() {
			h.worker()
		}, time.Second, ctx.Done())
	}

	<-ctx.Done()
	return nil
}

func (h *Handler) syncHandler(event event) error {
	key := event.key

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return xerrors.Errorf("failed to synchronize handler: %w", err)
	}

	obj, err := h.informer.Get(namespace, name)
	if err != nil {
		// Resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			h.ktx.Errorf("%s '%s' in work queue no longer exists: %s", strings.ToLower(h.kind.Name()), key, err)
			return nil
		}
		return xerrors.Errorf("failed to synchronize handler: %w", err)
	}
	objId := fmt.Sprintf("%s/%s@%s", obj.GetNamespace(), obj.GetName(), obj.GetResourceVersion())

	var handler handlerFunc
	switch event._type {
	case createEvent:
		handler = handlerFunc(h.events.CreateHandlerFunc)
	case updateEvent:
		handler = handlerFunc(h.events.UpdateHandlerFunc)
	case deleteEvent:
		handler = handlerFunc(h.events.DeleteHandlerFunc)
	default:
		h.ktx.Warnf("Invalid event type '%i' for %s... event skipped", event._type, objId)
		return nil
	}

	if handler == nil {
		return nil
	}

	if err = handler(h.ktx.SubContext(key), obj); err != nil {
		return err
	}

	if robj, ok := obj.(runtime.Object); ok {
		h.recorder.Eventf(
			robj,
			corev1.EventTypeNormal, "Synced",
			"Successfully synced by %s", h.ktx.Value(KontextKey("name")).(string),
		)
	}
	return nil
}

func (h *Handler) handleErr(err error, key interface{}) {
	if err == nil {
		h.queue.Forget(key)
		return
	}

	if h.queue.NumRequeues(key) < 10 {
		h.queue.AddRateLimited(key)
		return
	}

	// If number of requeue is above maxRetries, we drop the element out the queue
	//utilruntime.HandleError(err)
	h.queue.Forget(key)
}
func (h *Handler) processNextWorkItem() bool {
	key, shutdown := h.queue.Get()
	if shutdown {
		return false
	}

	defer h.queue.Done(key)
	err := h.syncHandler(key.(event))
	// TODO: NumRequestLimit must be defined by the user
	if err != nil && h.queue.NumRequeues(key) < 10 {
		h.queue.AddRateLimited(key)
		return true
	}

	h.queue.Forget(key)
	return true
}
func (h *Handler) worker() {
	for h.processNextWorkItem() {
	}
}
