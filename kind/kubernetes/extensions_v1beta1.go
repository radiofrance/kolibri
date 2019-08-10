package kubernetes

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/informers/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/radiofrance/kolibri/kind"
)

type ExtensionsV1Beta1 struct{ kubernetesAPI }

func (ExtensionsV1Beta1) APIVersion() string { return "extensions/v1beta1" }

type ingress struct{ ExtensionsV1Beta1 }
type ingressInformer struct {
	informer v1beta1.IngressInformer
	factory  informers.SharedInformerFactory
}

func (e ExtensionsV1Beta1) Ingress() *ingress { return &ingress{ExtensionsV1Beta1: e} }

func (ingress) Name() string { return "Ingress" }
func (k ingress) Informer(resync time.Duration, options ...informers.SharedInformerOption) kind.Informer {
	factory := informers.NewSharedInformerFactoryWithOptions(k.client.(kubernetes.Interface), resync, options...)
	return &ingressInformer{informer: factory.Extensions().V1beta1().Ingresses(), factory: factory}
}
func (i ingressInformer) Informer() interface{} { return i.informer }
func (i ingressInformer) HasSynced() bool       { return i.informer.Informer().HasSynced() }
func (i ingressInformer) Get(namespace, name string) (metav1.Object, error) {
	return i.informer.Lister().Ingresses(namespace).Get(name)
}
func (i ingressInformer) AddEventHandler(handler cache.ResourceEventHandler) {
	i.informer.Informer().AddEventHandler(handler)
}
func (i ingressInformer) Start(chanStop <-chan struct{}) { i.factory.Start(chanStop) }
