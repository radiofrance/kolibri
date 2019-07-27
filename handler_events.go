package kolibri

import (
	"golang.org/x/xerrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

type UpdateHandlerPolicy func(old, new metav1.Object) bool
type EventOption func(ktr *Kontroller, events *cache.ResourceEventHandlerFuncs)

// CreateHandlerFunc is a function that handle Kubernetes resources creation.
type CreateHandlerFunc func(ktx *Kontext, curr metav1.Object) error

func OnCreate(fnc CreateHandlerFunc) EventOption {
	return func(ktr *Kontroller, events *cache.ResourceEventHandlerFuncs) {
		events.AddFunc = func(obj interface{}) {
			ktx := ktr.newContext("OnCreate")

			curr, err := baseHandler(ktx, obj)
			if err != nil {
				ktr.handleError(err)
				return
			}

			err = fnc(ktx, curr)
			if err != nil {
				ktr.handleError(err)
			}
		}
	}
}

// UpdateHandlerFunc is a function that handle Kubernetes resources update.
type UpdateHandlerFunc func(ktx *Kontext, old, curr metav1.Object) error

func OnChange(fnc UpdateHandlerFunc) EventOption {
	return func(ktr *Kontroller, events *cache.ResourceEventHandlerFuncs) {
		events.UpdateFunc = func(oldObj interface{}, currObj interface{}) {
			ktx := ktr.newContext("OnChange")

			old, err := baseHandler(ktx, oldObj)
			if err != nil {
				ktr.handleError(err)
				return
			}
			curr, err := baseHandler(ktx, currObj)
			if err != nil {
				ktr.handleError(err)
				return
			}

			if !ktr.updatePolicy(old, curr) {
				ktr.Debugf("ignore '%s/%s' update due to the update policy", old.GetNamespace(), old.GetName())
				return
			}

			err = fnc(ktx, old, curr)
			if err != nil {
				ktr.handleError(err)
			}
		}
	}
}

// DeleteHandlerFunc is a function that handle Kubernetes resources deletion.
type DeleteHandlerFunc func(ktx *Kontext, curr metav1.Object) error

func OnDelete(fnc DeleteHandlerFunc) EventOption {
	return func(ktr *Kontroller, events *cache.ResourceEventHandlerFuncs) {
		events.DeleteFunc = func(obj interface{}) {
			ktx := ktr.newContext("OnDelete")

			curr, err := baseHandler(ktx, obj)
			if err != nil {
				ktr.handleError(err)
				return
			}

			err = fnc(ktx, curr)
			if err != nil {
				ktr.handleError(err)
			}
		}
	}
}

func baseHandler(ktx *Kontext, obj interface{}) (metav1.Object, error) {
	var object metav1.Object
	var ok bool

	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return nil, xerrors.Errorf("error decoding object, invalid type")
		}
		ktx.Debugf("tombstone found", object.GetName())
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			return nil, xerrors.Errorf("error decoding object tombstone, invalid type")
		}
		ktx.Debugf("recovered deleted object '%s' from tombstone", object.GetName())
	}
	return object, nil
}
