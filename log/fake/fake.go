package fake

import (
	"github.com/radiofrance/kolibri/log"
)

type Fake struct{ string }

func New() log.Logger { return &Fake{} }

func (f Fake) With(...log.Field) log.Logger { return &Fake{f.string} }
func (f Fake) Named(s string) log.Logger    { return &Fake{s} }

func (Fake) Debug(string)                  {}
func (Fake) Debugf(string, ...interface{}) {}
func (Fake) Info(string)                   {}
func (Fake) Infof(string, ...interface{})  {}
func (Fake) Warn(string)                   {}
func (Fake) Warnf(string, ...interface{})  {}
func (Fake) Error(string)                  {}
func (Fake) Errorf(string, ...interface{}) {}
func (Fake) Fatal(string)                  {}
func (Fake) Fatalf(string, ...interface{}) {}
func (Fake) Panic(string)                  {}
func (Fake) Panicf(string, ...interface{}) {}
