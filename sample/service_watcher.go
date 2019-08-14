package main

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/radiofrance/kolibri"
	kkind "github.com/radiofrance/kolibri/kind/kubernetes"
	"github.com/radiofrance/kolibri/log"
	"github.com/radiofrance/kolibri/log/klogrus"
)

func handler(ktx *kolibri.Kontext, event string) error {
	obj, err := ktx.Object()
	if err != nil {
		if errors.IsNotFound(err) {
			ktx.Infof("Event handled: %s ... but no object found (%s)", event, err)
			return nil
		}
		ktx.Errorf("Failed to get kubernetes object: %s", err)
		return err
	}

	ktx.
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
		config, err = clientcmd.BuildConfigFromFlags("", "/home/xunleii/.kube/config")
		handleErr(err)
	}

	client, err := kubernetes.NewForConfig(config)
	handleErr(err)

	ktr, err := kolibri.NewController("service_watcher", client)
	handleErr(err)

	ktr.SetLogger(klogrus.New(logrus.New()))

	svc, err := ktr.NewHandler(
		kkind.CoreV1(client).Service(),
		kolibri.OnAllNamespaces(),

		kolibri.WithUpdatePolicy(func(old, new v1.Object) bool { return old.GetResourceVersion() != new.GetResourceVersion() }),

		kolibri.OnCreate(func(ktx *kolibri.Kontext) error { return handler(ktx, "ServiceCreation") }),
		kolibri.OnChange(func(ktx *kolibri.Kontext) error { return handler(ktx, "ServiceUpdate") }),
		kolibri.OnDelete(func(ktx *kolibri.Kontext) error { return handler(ktx, "ServiceDeletion") }),
	)
	handleErr(err)

	err = ktr.Register(svc)
	handleErr(err)

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Minute)
	err = ktr.Run(ctx)
	handleErr(err)
}
