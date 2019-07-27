package kzap

import (
	"fmt"

	"go.uber.org/zap"

	"gitlab.com/xunleii.io/kolibri/log"
)

type kzap struct {
	core *zap.Logger
}

func New(logger *zap.Logger) log.Logger {
	return &kzap{logger}
}

func (k kzap) With(fields ...log.Field) log.Logger {
	zfields := make([]zap.Field, len(fields))

	for i, field := range fields {
		zfields[i] = zap.Field(field)
	}
	new := k.core.With(zfields...)
	return &kzap{new}
}

func (k kzap) Named(s string) log.Logger {
	new := k.core.Named(s)
	return &kzap{new}
}

func (k kzap) Debug(message string)                   { k.core.Debug(message) }
func (k kzap) Debugf(format string, a ...interface{}) { k.Debug(fmt.Sprintf(format, a...)) }

func (k kzap) Info(message string)                   { k.core.Info(message) }
func (k kzap) Infof(format string, a ...interface{}) { k.Info(fmt.Sprintf(format, a...)) }

func (k kzap) Warn(message string)                   { k.core.Warn(message) }
func (k kzap) Warnf(format string, a ...interface{}) { k.Warn(fmt.Sprintf(format, a...)) }

func (k kzap) Error(message string)                   { k.core.Error(message) }
func (k kzap) Errorf(format string, a ...interface{}) { k.Error(fmt.Sprintf(format, a...)) }

func (k kzap) Fatal(message string)                   { k.core.Fatal(message) }
func (k kzap) Fatalf(format string, a ...interface{}) { k.Fatal(fmt.Sprintf(format, a...)) }

func (k kzap) Panic(message string)                   { k.core.Panic(message) }
func (k kzap) Panicf(format string, a ...interface{}) { k.Panic(fmt.Sprintf(format, a...)) }
