package kolibri

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/radiofrance/kolibri/kind"
	"github.com/radiofrance/kolibri/log"
)

type Kontext struct {
	context.Context
	log.Logger

	informer  kind.Informer
	namespace string
	name      string
}

type KontextKey string

func (ktx *Kontext) SubContext(name string) *Kontext {
	return &Kontext{
		Context:  ktx.Context,
		Logger:   ktx.Named(name),
		informer: ktx.informer,
	}
}

func (ktx *Kontext) HandlerContext(namespace, name string) *Kontext {
	return &Kontext{
		Context:   ktx.Context,
		Logger:    ktx.Named(fmt.Sprintf("%s/%s", namespace, name)),
		informer:  ktx.informer,
		namespace: namespace,
		name:      name,
	}
}

func (ktx *Kontext) Namespace() string          { return ktx.namespace }
func (ktx *Kontext) Name() string               { return ktx.name }
func (ktx *Kontext) Object() (v1.Object, error) { return ktx.informer.Get(ktx.namespace, ktx.name) }
