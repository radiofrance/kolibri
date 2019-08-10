package kolibri

import (
	"context"

	"github.com/radiofrance/kolibri/log"
)

type Kontext struct {
	context.Context
	log.Logger
}

type KontextKey string

func (ktx *Kontext) SubContext(name string) *Kontext {
	return &Kontext{
		Context: ktx.Context,
		Logger:  ktx.Named(name),
	}
}
