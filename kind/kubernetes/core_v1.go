package kubernetes

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	corev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/radiofrance/kolibri/kind"
)

type CoreV1Interface struct{ kubernetesAPI }

func (CoreV1Interface) APIVersion() string { return "v1" }

func CoreV1(client kubernetes.Interface) *CoreV1Interface {
	return &CoreV1Interface{kubernetesAPI{client: client}}
}

type service struct{ CoreV1Interface }
type serviceInformer struct {
	informer corev1.ServiceInformer
	factory  informers.SharedInformerFactory
}

func (c CoreV1Interface) Service() *service { return &service{CoreV1Interface: c} }

func (service) Name() string { return "Service" }
func (k service) Informer(resync time.Duration, options ...informers.SharedInformerOption) kind.Informer {
	factory := informers.NewSharedInformerFactoryWithOptions(k.client.(kubernetes.Interface), resync, options...)
	return &serviceInformer{informer: factory.Core().V1().Services(), factory: factory}
}
func (i serviceInformer) Informer() interface{} { return i.informer }
func (i serviceInformer) HasSynced() bool       { return i.informer.Informer().HasSynced() }
func (i serviceInformer) Get(namespace, name string) (metav1.Object, error) {
	return i.informer.Lister().Services(namespace).Get(name)
}
func (i serviceInformer) AddEventHandler(handler cache.ResourceEventHandler) {
	i.informer.Informer().AddEventHandler(handler)
}
func (i serviceInformer) Start(chanStop <-chan struct{}) { i.factory.Start(chanStop) }
