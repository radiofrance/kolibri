package kolibri

import (
	"golang.org/x/xerrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// eventRegistry contains user handler function to be called when
// an action occurs in kubernetes.
type eventRegistry struct {
	createFuncs []CreateHandlerFunc
	updateFuncs []UpdateHandlerFunc
	deleteFuncs []DeleteHandlerFunc
}

// eventOption wraps functions defining on to handle kubernetes event on the watched object.
type eventOption func(events *eventRegistry) error
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

// handlerFunc is a generic function that handle an event.
type handlerFunc func(ktx *Kontext, curr metav1.Object) error

// CreateHandlerFunc is a function that handle Kubernetes resources creation.
type CreateHandlerFunc handlerFunc

// OnCreate registers function which will be called each time a new watched
// object will be created.
func OnCreate(fnc CreateHandlerFunc) eventOption {
	return func(events *eventRegistry) error {
		if events.CreateHandlerFunc != nil {
			return xerrors.New("OnCreate can only be called once")
		}
		events.CreateHandlerFunc = fnc
		return nil
	}
}

// UpdateHandlerFunc is a function that handle Kubernetes resources update.
type UpdateHandlerFunc handlerFunc

// OnChange registers function which will be called each time a watched
// object will be updated and validated by the policy.
func OnChange(fnc UpdateHandlerFunc) eventOption {
	return func(events *eventRegistry) error {
		if events.UpdateHandlerFunc != nil {
			return xerrors.New("OnChange can only be called once")
		}
		events.UpdateHandlerFunc = fnc
		return nil
	}
}

// DeleteHandlerFunc is a function that handle Kubernetes resources deletion.
type DeleteHandlerFunc handlerFunc

// OnDelete registers function which will be called each time a watched
// object will be removed.
func OnDelete(fnc DeleteHandlerFunc) eventOption {
	return func(events *eventRegistry) error {
		if events.DeleteHandlerFunc != nil {
			return xerrors.New("OnDelete can only be called once")
		}
		events.DeleteHandlerFunc = fnc
		return nil
	}
}
