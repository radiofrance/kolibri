package kind

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type Kind interface {
	APIVersion() string
	Name() string

	Informer(resync time.Duration, options ...informers.SharedInformerOption) Informer
}

type Informer interface {
	AddEventHandler(handler cache.ResourceEventHandler)

	Informer() cache.SharedIndexInformer
	HasSynced() bool

	Get(namespace, name string) (metav1.Object, error)

	Start(<-chan struct{})
}

func FullName(kind Kind) string { return kind.APIVersion() + "/" + kind.Name() }
