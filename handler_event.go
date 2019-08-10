package kolibri

import (
	"golang.org/x/xerrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

type eventType byte
type event struct {
	key   string
	_type eventType
}

const (
	unknownEvent eventType = iota
	createEvent
	updateEvent
	deleteEvent
)

func (h *Handler) enqueueWith(_type eventType, object metav1.Object) error {
	key, err := cache.MetaNamespaceKeyFunc(object)
	if err != nil {
		return xerrors.Errorf("failed to enqueue '%s/%s@%s': %w", object.GetNamespace(), object.GetName(), object.GetResourceVersion(), err)
	}
	h.queue.Add(event{key, _type})
	return nil
}

func (h *Handler) baseHandler(obj interface{}) (metav1.Object, error) {
	var kobj metav1.Object
	var ok bool

	if kobj, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return nil, xerrors.New("error decoding object, invalid type")
		}
		h.ktx.Debugf("tombstone found", kobj.GetName())
		kobj, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			return nil, xerrors.New("error decoding object tombstone, invalid type")
		}
		h.ktx.Debugf("recovered deleted object '%s' from tombstone", kobj.GetName())
	}
	return kobj, nil
}

func (h *Handler) addHandler(curr interface{}) {
	currObj, err := h.baseHandler(curr)
	if err != nil {
		h.ktx.Errorf("failed to get Kubernetes object: %s", err)
		return
	}

	if err = h.enqueueWith(createEvent, currObj); err != nil {
		h.ktx.Errorf("failed to enqueue object: %s", err)
	}
}

func (h *Handler) updateHandler(old, new interface{}) {
	oldObj, err := h.baseHandler(old)
	if err != nil {
		h.ktx.Errorf("failed to get Kubernetes object: %s", err)
		return
	}

	newObj, err := h.baseHandler(new)
	if err != nil {
		h.ktx.Errorf("failed to get Kubernetes object: %s", err)
		return
	}

	if !h.updatePolicy(oldObj, newObj) {
		return
	}

	if err = h.enqueueWith(updateEvent, newObj); err != nil {
		h.ktx.Errorf("failed to enqueue object: %s", err)
	}
}

func (h *Handler) deleteHandler(curr interface{}) {
	currObj, err := h.baseHandler(curr)
	if err != nil {
		h.ktx.Errorf("failed to get Kubernetes object: %s", err)
		return
	}

	if err = h.enqueueWith(deleteEvent, currObj); err != nil {
		h.ktx.Errorf("failed to enqueue object: %s", err)
	}
}
