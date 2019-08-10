package fake

import (
	"path"

	"github.com/radiofrance/kolibri/log"
)

type Fake struct {
	string
	bool
}

func New() log.Logger { return &Fake{} }

func (f Fake) With(...log.Field) log.Logger { return &Fake{f.string, f.bool} }
func (f Fake) Named(s string) log.Logger    { return &Fake{path.Join(f.string, s), f.bool} }

func (Fake) Debug(string)                     {}
func (Fake) Debugf(string, ...interface{})    {}
func (Fake) Info(string)                      {}
func (Fake) Infof(string, ...interface{})     {}
func (Fake) Warn(string)                      {}
func (Fake) Warnf(string, ...interface{})     {}
func (Fake) Error(string)                     {}
func (Fake) Errorf(string, ...interface{})    {}
func (Fake) Fatal(string)                     {}
func (Fake) Fatalf(string, ...interface{})    {}
func (f *Fake) Panic(string)                  { f.bool = true }
func (f *Fake) Panicf(string, ...interface{}) { f.bool = true }
