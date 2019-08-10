package kolibri

import (
	"golang.org/x/xerrors"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/clientcmd"
)

// OnAllNamespaces configures the current handler to watch all namespaces (default behavior).
func OnAllNamespaces() Option { return OnNamespace("") }

// OnNamespace configures the current handler to watch only the specified namespace.
// Only one can be provided.
func OnNamespace(ns string) Option {
	return optionFnc(func(ctx *handlerBuildContext) error {
		ctx.informerOpts = append(ctx.informerOpts, informers.WithNamespace(ns))
		return nil
	})
}

// OnCurrentNamespace configures the current handler to watch the namespace on
// which the controller runs.
func OnCurrentNamespace(c clientcmd.ClientConfig) Option {
	return optionFnc(func(ctx *handlerBuildContext) error {
		if c == nil {
			return xerrors.Errorf("client config cannot be nil")
		}
		ns, _, err := c.Namespace()
		if err != nil {
			return err
		}

		ctx.informerOpts = append(ctx.informerOpts, informers.WithNamespace(ns))
		return nil
	})
}

// WithUpdatePolicy sets the update policy used by the controller to known when
// an object is considered as updated.
func WithUpdatePolicy(policy UpdateHandlerPolicy) Option {
	return optionFnc(func(ctx *handlerBuildContext) error {
		if policy == nil {
			return xerrors.Errorf("update policy cannot be nil")
		}

		ctx.ktrlOpts = append(ctx.ktrlOpts, func(ktr *Kontroller) error {
			ktr.setUpdatePolicy(policy)
			return nil
		})
		return nil
	})
}

// OnCreate registers function which will be called each time a new watched
// object will be created.
func OnCreate(fnc CreateHandlerFunc) Option {
	return optionFnc(func(ctx *handlerBuildContext) error {
		if fnc == nil {
			return xerrors.New("OnCreate function cannot be nil")
		}

		ctx.eventOpts = append(ctx.eventOpts, func(reg *eventRegistry) error {
			if reg.CreateHandlerFunc != nil {
				return xerrors.New("OnCreate can only be called once")
			}
			reg.CreateHandlerFunc = fnc
			return nil
		})
		return nil
	})
}

// OnChange registers function which will be called each time a watched
// object will be updated and validated by the policy.
func OnChange(fnc UpdateHandlerFunc) Option {
	return optionFnc(func(ctx *handlerBuildContext) error {
		if fnc == nil {
			return xerrors.New("OnChange function cannot be nil")
		}

		ctx.eventOpts = append(ctx.eventOpts, func(reg *eventRegistry) error {
			if reg.UpdateHandlerFunc != nil {
				return xerrors.New("OnChange can only be called once")
			}
			reg.UpdateHandlerFunc = fnc
			return nil
		})
		return nil
	})
}

// OnDelete registers function which will be called each time a watched
// object will be removed.
func OnDelete(fnc DeleteHandlerFunc) Option {
	return optionFnc(func(ctx *handlerBuildContext) error {
		if fnc == nil {
			return xerrors.New("OnDelete function cannot be nil")
		}

		ctx.eventOpts = append(ctx.eventOpts, func(reg *eventRegistry) error {
			if reg.DeleteHandlerFunc != nil {
				return xerrors.New("OnDelete can only be called once")
			}
			reg.DeleteHandlerFunc = fnc
			return nil
		})
		return nil
	})
}
