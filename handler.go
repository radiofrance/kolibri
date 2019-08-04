package kolibri

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/xerrors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	"github.com/radiofrance/kolibri/kind"
	"github.com/radiofrance/kolibri/log"
)

type Handler struct {
	ktr    *Kontroller
	events eventRegistry

	kind     kind.Kind
	informer kind.Informer

	queue    workqueue.RateLimitingInterface
	recorder record.EventRecorder
}

// handlerBuildContext contains all elements used to build an handler.
// Theses elements are provided by Option interface.
type handlerBuildContext struct {
	kind kind.Kind

	informerOpts []informers.SharedInformerOption
	ktrlOpts     kontrolerOptions
	eventOpts    eventOptions
}
type Option interface{ apply(ctx *handlerBuildContext) error }

func (k *Kontroller) NewHandler(opts ...Option) (*Handler, error) {
	ctx := &handlerBuildContext{}

	for _, opt := range opts {
		if err := opt.apply(ctx); err != nil {
			return nil, err
		}
	}

	kindOpt := kindOption(nil)
	factoryOpts := informerFactoryOptions{}
	ktrlOpts := kontrolerOptions{}
	eventOpts := eventOptions{}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case kindOption:
			if kindOpt != nil {
				return nil, xerrors.Errorf("only one kind must be provided")
			}
			kindOpt = opt
		case informerFactoryOption:
			factoryOpts = append(factoryOpts, opt)
		case kontrolerOption:
			ktrlOpts = append(ktrlOpts, opt)
		case eventOption:
			eventOpts = append(eventOpts, opt)
		default:
			return nil, xerrors.Errorf("%T not implemented", opt)
		}
	}

	if kindOpt == nil {
		return nil, xerrors.Errorf("kind must be provided")
	}
	kind, err := kindOpt()
	if err != nil {
		return nil, err
	}

	if len(eventOpts) == 0 {
		return nil, xerrors.Errorf("at least one event handler (On...) must be provided")
	}

	handler := &Handler{ktr: k.copy()}
	k = k.copy()
	k.Logger = k.Named(fmt.Sprintf("%s/%s", kind.APIVersion(), kind.Name()))

	client, err := k.client(kind.ClientType())
	if err != nil {
		return nil, err
	}

	informerOpts, err := factoryOpts.apply()
	if err != nil {
		return nil, err
	}
	handler.informer = kind.Informer(client, 5*time.Second, informerOpts...)

	err = ktrlOpts.apply(k)
	if err != nil {
		return nil, err
	}

	err = eventOpts.apply(&handler.events)
	if err != nil {
		return nil, err
	}

	handler.queue = workqueue.NewNamedRateLimitingQueue(
		workqueue.DefaultControllerRateLimiter(),
		fmt.Sprintf("%s:%s:%s/%s@%d", "kolibris", k.name, kind.APIVersion(), kind.Name(), uuid.New().String()),
	)

	enqueuWith := func(container eventContainer, object metav1.Object) {
		var key string
		var err error
		if key, err = cache.MetaNamespaceKeyFunc(object); err != nil {
			return
		}
		container.setKey(key)
		handler.queue.Add(container)
	}

	defaultUpdatePolicy := func(old, new metav1.Object) bool { return old.GetResourceVersion() != new.GetResourceVersion() }

	// -- Generic 'add' handler
	addHandler := func(obj interface{}, kind string) {
		object, err := baseHandler(handler.ktr.newContext("addHandler"), obj)
		if err != nil {
			return
		}
		enqueuWith(&createEvent{baseEvent: &baseEvent{kind: kind}}, object)
	}
	// -- Generic 'update' handler
	updateHandler := func(old, new interface{}, kind string) {
		oldObject, oerr := baseHandler(handler.ktr.newContext("updateHandler"), old)
		newObject, nerr := baseHandler(handler.ktr.newContext("updateHandler"), new)
		if oerr != nil || nerr != nil {
			return
		}
		if defaultUpdatePolicy(oldObject, newObject) {
			enqueuWith(&updateEvent{baseEvent: &baseEvent{kind: kind}}, newObject)
		}
	}
	// -- Generic 'delete' handler
	deleteHandler := func(obj interface{}, kind string) {
		object, err := baseHandler(handler.ktr.newContext("deleteHandler"), obj)
		if err != nil {
			return
		}
		enqueuWith(&deleteEvent{baseEvent: &baseEvent{kind: kind}}, object)
	}

	// -- Add event handler
	handler.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { addHandler(obj, kind.Name()) },
		UpdateFunc: func(old, new interface{}) { updateHandler(old, new, kind.Name()) },
		DeleteFunc: func(obj interface{}) { deleteHandler(obj, kind.Name()) },
	})

	kk := k.kube.(*kubernetes.Clientset)
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(k.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kk.CoreV1().Events("")})
	_ = eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "kolibri"})

	return handler, nil
}

func (h Handler) Run(ctx context.Context) error {
	defer h.queue.ShutDown()
	chanStop := make(chan struct{}, 10)

	h.informer.Start(chanStop)

	if ok := cache.WaitForCacheSync(chanStop, h.informer.HasSynced); !ok {
		panic("failed to wait for caches to sync")
	}

	for i := 0; i < 10; i++ {
		go wait.Until(func() {
			h.worker()
		}, time.Second, chanStop)
	}

	<-chanStop
	return nil
}

// FIXME clean theses function ... only for testing purpose
type eventContainer interface {
	Kind() string
	Key() string
	setKey(key string)
}
type createEvent struct{ *baseEvent }
type updateEvent struct{ *baseEvent }
type deleteEvent struct{ *baseEvent }

type baseEvent struct {
	kind string
	key  string
}

func (b baseEvent) Kind() string       { return b.kind }
func (b baseEvent) Key() string        { return b.key }
func (b *baseEvent) setKey(key string) { b.key = key }

// ---------------------------------------------------------------------------------------------------------------//
// Defines running controller processes
func (h *Handler) syncHandler(container eventContainer) error {
	key := container.Key()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	obj, err := h.informer.Get(namespace, name)
	if err != nil {
		// The Foo resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			// TODO: Handle error
			h.ktr.With(log.Error("err", err)).Errorf("service '%s' in work queue no longer exists", key)
			return nil
		}
		return err
	}

	var handler handlerFunc
	switch container.(type) {
	case *createEvent:
		handler = handlerFunc(h.events.CreateHandlerFunc)
	case *updateEvent:
		handler = handlerFunc(h.events.UpdateHandlerFunc)
	case *deleteEvent:
		handler = handlerFunc(h.events.DeleteHandlerFunc)
	}

	if handler == nil {
		return nil
	}
	err = handler(h.ktr.newContext(key), obj)
	if err != nil {
		return err
	}

	//eventRecorder.Event(svc, corev1.EventTypeNormal, "Synced", "Service synced successfully")
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
	event, shutdown := h.queue.Get()
	if shutdown {
		return false
	}

	defer h.queue.Done(event)
	err := h.syncHandler(event.(eventContainer))
	h.handleErr(err, event)

	return true
}
func (h *Handler) worker() {
	for h.processNextWorkItem() {
	}
}

// baseHandler is a generic handler used to extract valid object from the
// given interface.
func baseHandler(ktx *Kontext, obj interface{}) (metav1.Object, error) {
	var object metav1.Object
	var ok bool

	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return nil, xerrors.Errorf("error decoding object, invalid type")
		}
		ktx.Debugf("tombstone found", object.GetName())
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			return nil, xerrors.Errorf("error decoding object tombstone, invalid type")
		}
		ktx.Debugf("recovered deleted object '%s' from tombstone", object.GetName())
	}
	return object, nil
}
