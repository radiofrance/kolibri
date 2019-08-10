package kolibri

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
)

type Option interface{ apply(ctx *handlerBuildContext) error }

// CreateHandlerFunc is a function that handle Kubernetes resources creation.
type CreateHandlerFunc handlerFunc

// UpdateHandlerFunc is a function that handle Kubernetes resources update.
type UpdateHandlerFunc handlerFunc

// DeleteHandlerFunc is a function that handle Kubernetes resources deletion.
type DeleteHandlerFunc handlerFunc

// handlerFunc is a generic function that handle an event.
type handlerFunc func(ktx *Kontext, obj metav1.Object) error

// eventRegistry contains functions called when
// an event occurs on kubernetes.
type eventRegistry struct {
	CreateHandlerFunc
	UpdateHandlerFunc
	DeleteHandlerFunc
}

// handlerBuildContext contains all elements used to build an handler.
// Theses elements are provided by Option interface.
type handlerBuildContext struct {
	informerOpts sharedInformerOptions
	ktrlOpts     kontrolerOptions
	eventOpts    eventOptions
}

type sharedInformerOptions []informers.SharedInformerOption

// kontrolerOption wraps functions used to configure the handler controller.
type kontrolerOption func(*Kontroller) error
type kontrolerOptions []kontrolerOption

func (o kontrolerOptions) apply(k *Kontroller) error {
	for _, opt := range o {
		err := opt(k)
		if err != nil {
			return err
		}
	}
	return nil
}

// eventOption wraps functions defining on to handle kubernetes event on the watched object.
type eventOption func(*eventRegistry) error
type eventOptions []eventOption

func (o eventOptions) apply(events *eventRegistry) error {
	for _, opt := range o {
		err := opt(events)
		if err != nil {
			return err
		}
	}
	return nil
}

// optionFnc implements Option interface around the apply function
type optionFnc func(ctx *handlerBuildContext) error

func (o optionFnc) apply(ctx *handlerBuildContext) error { return o(ctx) }
