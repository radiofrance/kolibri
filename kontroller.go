package kolibri

import (
	"context"

	"golang.org/x/sync/errgroup"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/radiofrance/kolibri/log"
	"github.com/radiofrance/kolibri/log/fake"
)

type Kontroller struct {
	name string

	log.Logger
	kube            kubernetes.Interface
	handlers        []*Handler
	updatePolicyFnc UpdateHandlerPolicy
}

func NewController(name string, kube kubernetes.Interface, opts ...interface{}) *Kontroller {
	return &Kontroller{
		name:   name,
		kube:   kube,
		Logger: fake.New(),
	}
}

// UpdateHandlerPolicy that defines when two kubernetes objects are different.
type UpdateHandlerPolicy func(old, new metav1.Object) bool

func (k *Kontroller) SetLogger(logger log.Logger) { k.Logger = logger }
func (k *Kontroller) Register(handlers ...*Handler) error {
	k.handlers = append(k.handlers, handlers...)
	return nil
}
func (k *Kontroller) Run(ctx context.Context) error {
	errg, ctx := errgroup.WithContext(ctx)

	for _, handler := range k.handlers {
		errg.Go(func() error { return handler.Run(ctx) })
	}

	return errg.Wait()
}

func (k *Kontroller) newContext(name string) *Kontext { return &Kontext{k} }
func (k *Kontroller) handleError(err error)           {}

//TODO: Make a real copy
func (k *Kontroller) copy() *Kontroller { return k }

func (k *Kontroller) setUpdatePolicy(policy UpdateHandlerPolicy)              { k.updatePolicyFnc = policy }
func (k *Kontroller) updatePolicy(old metav1.Object, curr metav1.Object) bool { return false }
