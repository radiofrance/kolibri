package kolibri

import (
	"context"
	"reflect"

	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"gitlab.com/xunleii.io/kolibri/log"
	"gitlab.com/xunleii.io/kolibri/log/fake"
)

type Kontroller struct {
	name string

	log.Logger
	kube     kubernetes.Interface
	handlers []*Handler
}

func NewController(name string, client kubernetes.Interface, opts ...interface{}) *Kontroller {
	return &Kontroller{
		name:   name,
		Logger: fake.New(),
		kube:   client,
	}
}

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

func (k *Kontroller) copy() *Kontroller { return k }

func (k *Kontroller) setUpdatePolicy(policy UpdateHandlerPolicy)      { return }
func (k *Kontroller) updatePolicy(old v1.Object, curr v1.Object) bool { return false }

func (k *Kontroller) client(clientType reflect.Type) (interface{}, error) { return k.kube, nil }
