package kolibri

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1"

	"gitlab.com/xunleii.io/kolibri/log"
	"gitlab.com/xunleii.io/kolibri/log/fake"
)

type Kontroller struct {
	log.Logger
}

func NewController(name string, opts ...interface{}) *Kontroller {
	return &Kontroller{
		Logger: fake.New(),
	}
}

func (k *Kontroller) SetLogger(logger log.Logger)        { panic("implement me") }
func (k *Kontroller) Register(handlers ...Handler) error { panic("implement me") }
func (k *Kontroller) Run(context context.Context) error  { panic("implement me") }

func (k *Kontroller) newContext(name string) *Kontext { panic("implement me") }
func (k *Kontroller) handleError(err error)           { panic("implement me") }

func (k *Kontroller) updatePolicy(old v1.Object, curr v1.Object) bool { return false }
