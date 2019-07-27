package kolibri

import (
	"context"

	"gitlab.com/xunleii.io/kolibri/log"
)

type Kontroller struct{}

func (k *Kontroller) SetLogger(logger log.Logger)        { panic("implement me") }
func (k *Kontroller) Register(handlers ...Handler) error { panic("implement me") }
func (k *Kontroller) Run(context context.Context) error  { panic("implement me") }
