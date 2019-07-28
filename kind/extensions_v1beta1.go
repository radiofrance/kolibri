package kind

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/informers/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
)

type Ingress struct{ *kubeKind }
type IngressInformer struct{ informer v1beta1.IngressInformer }

func (Ingress) APIVersion() string { return "extensions/v1beta1" }
func (Ingress) Name() string       { return "Ingress" }
func (Ingress) Informer(client interface{}, resync time.Duration, options ...informers.SharedInformerOption) Informer {
	factory := informers.NewSharedInformerFactoryWithOptions(client.(kubernetes.Interface), resync, options...)
	return &IngressInformer{informer: factory.Extensions().V1beta1().Ingresses()}
}
func (i IngressInformer) Informer() interface{} { return i.informer }
func (i IngressInformer) HasSynced() bool       { return i.informer.Informer().HasSynced() }
func (i IngressInformer) Get(namespace, name string) (v1.Object, error) {
	return i.informer.Lister().Ingresses(namespace).Get(name)
}
