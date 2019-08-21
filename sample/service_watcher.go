package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/radiofrance/kolibri"
	kind "github.com/radiofrance/kolibri/kind/kubernetes"
	"github.com/radiofrance/kolibri/log"
	"github.com/radiofrance/kolibri/log/klogrus"
)

var kubeconfig *string

func init() {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
}

func main() {
	// Configure kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		handleErr(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	handleErr(err)

	// Create new Kolibri Controller
	ktr, err := kolibri.NewController("service_watcher", clientset)
	handleErr(err)

	ktr.SetLogger(klogrus.New(logrus.New()))

	// Create the service handler
	svc, err := ktr.NewHandler(
		kind.CoreV1(clientset).Service(),
		kolibri.OnAllNamespaces(),

		kolibri.WithUpdatePolicy(kolibri.DefaultUpdateHandlerPolicy),

		kolibri.OnCreate(func(ktx *kolibri.Kontext) error { return handler(ktx, "ServiceCreation") }),
		kolibri.OnChange(func(ktx *kolibri.Kontext) error { return handler(ktx, "ServiceUpdate") }),
		kolibri.OnDelete(func(ktx *kolibri.Kontext) error { return handler(ktx, "ServiceDeletion") }),
	)
	handleErr(err)

	// Run the controller
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	err = ktr.Run(ctx, svc)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

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
