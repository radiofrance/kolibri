package kind

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	corev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type corev1Kind struct{ kubeKind }

func (corev1Kind) APIVersion() string { return "v1" }

type Service struct{ corev1Kind }
type ServiceInformer struct {
	informer corev1.ServiceInformer
	factory  informers.SharedInformerFactory
}

func (Service) Name() string { return "Service" }
func (Service) Informer(client interface{}, resync time.Duration, options ...informers.SharedInformerOption) Informer {
	factory := informers.NewSharedInformerFactoryWithOptions(client.(kubernetes.Interface), resync, options...)
	return &ServiceInformer{informer: factory.Core().V1().Services(), factory: factory}
}
func (i ServiceInformer) Informer() interface{} { return i.informer }
func (i ServiceInformer) HasSynced() bool       { return i.informer.Informer().HasSynced() }
func (i ServiceInformer) Get(namespace, name string) (metav1.Object, error) {
	return i.informer.Lister().Services(namespace).Get(name)
}
func (i ServiceInformer) AddEventHandler(handler cache.ResourceEventHandler) {
	i.informer.Informer().AddEventHandler(handler)
}
func (i ServiceInformer) Start(chanStop <-chan struct{}) { i.factory.Start(chanStop) }
