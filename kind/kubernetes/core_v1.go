package kubernetes

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/radiofrance/kolibri/kind"
)

type CoreV1Interface struct{ kubernetesAPI }

func (CoreV1Interface) APIVersion() string { return "corev1" }

func CoreV1(client kubernetes.Interface) *CoreV1Interface {
	return &CoreV1Interface{kubernetesAPI{client: client}}
}

type service struct{ CoreV1Interface }
type serviceInformer struct {
	factory  informers.SharedInformerFactory
	informer cache.SharedIndexInformer
	lister   corev1.ServiceLister
}

func (c CoreV1Interface) Service() *service { return &service{CoreV1Interface: c} }

func (service) Name() string { return "Service" }
func (k service) Informer(resync time.Duration, options ...informers.SharedInformerOption) kind.Informer {
	factory := informers.NewSharedInformerFactoryWithOptions(k.client.(kubernetes.Interface), resync, options...)
	svc := factory.Core().V1().Services()
	return &serviceInformer{factory: factory, informer: svc.Informer(), lister: svc.Lister()}
}
func (i serviceInformer) Informer() cache.SharedIndexInformer { return i.informer }
func (i serviceInformer) HasSynced() bool                     { return i.informer.HasSynced() }
func (i serviceInformer) Get(namespace, name string) (metav1.Object, error) {
	return i.lister.Services(namespace).Get(name)
}
func (i serviceInformer) AddEventHandler(handler cache.ResourceEventHandler) {
	i.informer.AddEventHandler(handler)
}
func (i serviceInformer) Start(chanStop <-chan struct{}) { i.factory.Start(chanStop) }

type pod struct{ CoreV1Interface }
type podInformer struct {
	factory  informers.SharedInformerFactory
	informer cache.SharedIndexInformer
	lister   corev1.PodLister
}

func (c CoreV1Interface) Pod() *pod { return &pod{CoreV1Interface: c} }

func (pod) Name() string { return "Pod" }
func (k pod) Informer(resync time.Duration, options ...informers.SharedInformerOption) kind.Informer {
	factory := informers.NewSharedInformerFactoryWithOptions(k.client.(kubernetes.Interface), resync, options...)
	pod := factory.Core().V1().Pods()
	return &podInformer{factory: factory, informer: pod.Informer(), lister: pod.Lister()}
}
func (i podInformer) Informer() cache.SharedIndexInformer { return i.informer }
func (i podInformer) HasSynced() bool                     { return i.informer.HasSynced() }
func (i podInformer) Get(namespace, name string) (metav1.Object, error) {
	return i.lister.Pods(namespace).Get(name)
}
func (i podInformer) AddEventHandler(handler cache.ResourceEventHandler) {
	i.informer.AddEventHandler(handler)
}
func (i podInformer) Start(chanStop <-chan struct{}) { i.factory.Start(chanStop) }
