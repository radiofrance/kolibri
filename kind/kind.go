package kind

import (
	"reflect"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
)

type Kind interface {
	ClientType() reflect.Type
	APIVersion() string
	Name() string

	Informer(client interface{}, resync time.Duration, options ...informers.SharedInformerOption) Informer
}

type Informer interface {
	Informer() interface{}
	HasSynced() bool
	Get(namespace, name string) (metav1.Object, error)
}
