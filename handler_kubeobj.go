package kolibri

import (
	"golang.org/x/xerrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/clientcmd"

	"gitlab.com/xunleii.io/kolibri/kind"
)

// kindOption wraps a function which verify the validity of a kind.
type kindOption func() (kind.Kind, error)

func (k kindOption) apply(ctx *handlerBuildContext) error {  }

// Kind register the 'kind' of the kubernetes object which we want to
// control.
// Only one must be registered.
func Kind(k kind.Kind) kindOption {
	return func() (i kind.Kind, e error) {
		if k == nil {
			return nil, xerrors.Errorf("kind cannot be nil")
		}
		return k, nil
	}
}

// informerFactoryOption wraps functions used to configure the InformerFactory.
type informerFactoryOption func() ([]informers.SharedInformerOption, error)
type informerFactoryOptions []informerFactoryOption

func (o informerFactoryOptions) apply() ([]informers.SharedInformerOption, error) {
	var opts []informers.SharedInformerOption

	for _, factOpt := range o {
		opt, err := factOpt()
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt...)
	}

	return opts, nil
}

// OnAllNamespaces configures the current handler to watch all namespaces (default behavior).
func OnAllNamespaces() informerFactoryOption {
	return func() ([]informers.SharedInformerOption, error) {
		return []informers.SharedInformerOption{
			informers.WithNamespace(""),
		}, nil
	}
}

// OnNamespace configures the current handler to watch only the specified namespace.
// Only one can be provided.
func OnNamespace(ns string) informerFactoryOption {
	return func() ([]informers.SharedInformerOption, error) {
		return []informers.SharedInformerOption{
			informers.WithNamespace(ns),
		}, nil
	}
}

// OnCurrentNamespace configures the current handler to watch the namespace on
// which the controller runs.
func OnCurrentNamespace(c clientcmd.ClientConfig) informerFactoryOption {
	return func() ([]informers.SharedInformerOption, error) {
		if c == nil {
			return nil, xerrors.Errorf("client config cannot be nil")
		}
		ns, _, err := c.Namespace()
		if err != nil {
			return nil, err
		}

		return []informers.SharedInformerOption{
			informers.WithNamespace(ns),
		}, nil
	}
}

// UpdateHandlerPolicy that defines when two kubernetes objects are different.
type UpdateHandlerPolicy func(old, new metav1.Object) bool

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

// WithUpdatePolicy sets the update policy used by the controller to known when
// an object is considered as updated.
func WithUpdatePolicy(policy UpdateHandlerPolicy) kontrolerOption {
	return func(ktr *Kontroller) error {
		if policy == nil {
			return xerrors.Errorf("update policy cannot be nil")
		}

		ktr.setUpdatePolicy(policy)
		return nil
	}
}
