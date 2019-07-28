package kind

import (
	"reflect"

	"k8s.io/client-go/kubernetes"
)

type kubeKind struct{}

func (kubeKind) ClientType() reflect.Type { return reflect.TypeOf(kubernetes.Interface(nil)) }
