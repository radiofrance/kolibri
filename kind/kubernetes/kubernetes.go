package kubernetes

import (
	"k8s.io/client-go/kubernetes"
)

type kubernetesAPI struct{ client kubernetes.Interface }
