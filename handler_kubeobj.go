package kolibri

import (
	"golang.org/x/xerrors"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/clientcmd"

	"gitlab.com/xunleii.io/kolibri/kind"
)

type KindOption func() (kind.Kind, error)

func Kind(k kind.Kind) KindOption {
	return func() (i kind.Kind, e error) {
		if k == nil {
			return nil, xerrors.Errorf("kind cannot be nil")
		}
		return k, nil
	}
}

type InformerFactoryOption func() ([]informers.SharedInformerOption, error)

func OnAllNamespaces() InformerFactoryOption {
	return func() ([]informers.SharedInformerOption, error) {
		return []informers.SharedInformerOption{
			informers.WithNamespace(""),
		}, nil
	}
}
func OnNamespace(ns string) InformerFactoryOption {
	return func() ([]informers.SharedInformerOption, error) {
		return []informers.SharedInformerOption{
			informers.WithNamespace(ns),
		}, nil
	}
}
func OnCurrentNamespace(c clientcmd.ClientConfig) InformerFactoryOption {
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
