package main

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/radiofrance/kolibri"
	"github.com/radiofrance/kolibri/kind"
	"github.com/radiofrance/kolibri/log"
	"github.com/radiofrance/kolibri/log/kzap"
)

func handler(ktx *kolibri.Kontext, event string, obj v1.Object) error {
	log.Logger(ktx).
		With(log.String("svc", fmt.Sprintf("%s/%s@%s (%s)", obj.GetNamespace(), obj.GetName(), obj.GetResourceVersion(), obj.GetUID()))).
		Infof("Event handled: %s", event)
	return nil
}

func main() {
	handleErr := func(err error) {
		if err != nil {
			panic(err.Error())
		}
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", "/home/anicolaie/.kube/config")
		handleErr(err)
	}

	client, err := kubernetes.NewForConfig(config)
	handleErr(err)

	ktr := kolibri.NewController("service_watcher", client)
	ktr.Logger = kzap.New(zap.NewExample())

	svc, err := ktr.NewHandler(
		kolibri.OnAllNamespaces(),
		kolibri.Kind(&kind.Service{}),

		kolibri.WithUpdatePolicy(func(old, new v1.Object) bool { return old.GetResourceVersion() != new.GetResourceVersion() }),

		kolibri.OnCreate(func(ktx *kolibri.Kontext, obj v1.Object) error { return handler(ktx, "ServiceCreation", obj) }),
		kolibri.OnChange(func(ktx *kolibri.Kontext, obj v1.Object) error { return handler(ktx, "ServiceUpdate", obj) }),
		kolibri.OnDelete(func(ktx *kolibri.Kontext, obj v1.Object) error { return handler(ktx, "ServiceDeletion", obj) }),
	)
	handleErr(err)

	err = ktr.Register(svc)
	handleErr(err)

	err = ktr.Run(context.TODO())
	handleErr(err)
}
