package main

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func main() {
	ctrl.SetLogger(zap.New(func(options *zap.Options) {
		options.Development = true
	}))

	// Create manager
	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		ctrl.Log.Error(err, "oops")
		return
	}

	// Create service reconciler
	err = ctrl.NewControllerManagedBy(manager).
		Named("ServiceReconcilier").
		WithEventFilter(LabelExistsFilter("watch-me-please")).
		For(&corev1.Service{}).
		Complete(ServiceWatcherReconciler{manager.GetClient()})
	if err != nil {
		ctrl.Log.Error(err, "oops")
		return
	}

	manager.Start(make(chan struct{}))
}

type LabelExistsFilter string

func (s LabelExistsFilter) Create(event event.CreateEvent) bool {
	_, exists := event.Meta.GetLabels()[string(s)]
	return exists
}

func (s LabelExistsFilter) Delete(event event.DeleteEvent) bool {
	_, exists := event.Meta.GetLabels()[string(s)]
	return exists
}

func (s LabelExistsFilter) Update(event event.UpdateEvent) bool {
	_, exists := event.MetaNew.GetLabels()[string(s)]
	return exists
}
func (s LabelExistsFilter) Generic(event event.GenericEvent) bool {
	_, exists := event.Meta.GetLabels()[string(s)]
	return exists
}

type ServiceWatcherReconciler struct{ client.Client }

func (r ServiceWatcherReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	log := ctrl.Log.WithValues("service", req.NamespacedName)

	svc := &corev1.Service{}
	err := r.Client.Get(ctx, req.NamespacedName, svc)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	log.Info("service found", "service", svc)

	if svc.Annotations == nil {
		svc.Annotations = map[string]string{}
	}
	svc.Annotations["service.example/watched"] = "true"
	err = r.Client.Update(ctx, svc)
	if err != nil {
		log.Error(err, "failed to update service")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
