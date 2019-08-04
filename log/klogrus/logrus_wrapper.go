package klogrus

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"

	"github.com/radiofrance/kolibri/log"
)

type core logrus.FieldLogger
type klogrus struct {
	core
	name string
}

var converter = map[zapcore.FieldType]func(logrus.Fields, log.Field){
	zapcore.UnknownType:         func(logrus.Fields, log.Field) {},
	zapcore.ArrayMarshalerType:  func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Interface },
	zapcore.ObjectMarshalerType: func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Interface },
	zapcore.BinaryType:          func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Interface.([]byte) },
	zapcore.BoolType:            func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Integer == 1 },
	zapcore.ByteStringType:      func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Interface.([]byte) },
	zapcore.Complex128Type:      func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Interface.(complex128) },
	zapcore.Complex64Type:       func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Interface.(complex64) },
	zapcore.DurationType:        func(lf logrus.Fields, f log.Field) { lf[f.Key] = time.Duration(f.Integer) },
	zapcore.Float64Type:         func(lf logrus.Fields, f log.Field) { lf[f.Key] = math.Float64frombits(uint64(f.Integer)) },
	zapcore.Float32Type:         func(lf logrus.Fields, f log.Field) { lf[f.Key] = math.Float32frombits(uint32(f.Integer)) },
	zapcore.Int64Type:           func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Integer },
	zapcore.Int32Type:           func(lf logrus.Fields, f log.Field) { lf[f.Key] = int32(f.Integer) },
	zapcore.Int16Type:           func(lf logrus.Fields, f log.Field) { lf[f.Key] = int16(f.Integer) },
	zapcore.Int8Type:            func(lf logrus.Fields, f log.Field) { lf[f.Key] = int8(f.Integer) },
	zapcore.StringType:          func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.String },
	zapcore.Uint64Type:          func(lf logrus.Fields, f log.Field) { lf[f.Key] = uint64(f.Integer) },
	zapcore.Uint32Type:          func(lf logrus.Fields, f log.Field) { lf[f.Key] = uint32(f.Integer) },
	zapcore.Uint16Type:          func(lf logrus.Fields, f log.Field) { lf[f.Key] = uint16(f.Integer) },
	zapcore.Uint8Type:           func(lf logrus.Fields, f log.Field) { lf[f.Key] = uint8(f.Integer) },
	zapcore.UintptrType:         func(lf logrus.Fields, f log.Field) { lf[f.Key] = uintptr(f.Integer) },
	zapcore.ReflectType:         func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Interface },
	zapcore.NamespaceType:       func(lf logrus.Fields, f log.Field) { panic("not implemented") },
	zapcore.StringerType:        func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Interface.(fmt.Stringer).String() },
	zapcore.ErrorType:           func(lf logrus.Fields, f log.Field) { lf[f.Key] = f.Interface },
	zapcore.SkipType:            func(logrus.Fields, log.Field) {},

	zapcore.TimeType: func(lf logrus.Fields, f log.Field) {
		if f.Interface != nil {
			lf[f.Key] = time.Unix(0, f.Integer).In(f.Interface.(*time.Location))
		} else {
			// Fall back to UTC if location is nil.
			lf[f.Key] = time.Unix(0, f.Integer)
		}
	},
}

func New(logger *logrus.Logger) log.Logger {
	return &klogrus{core: logger}
}

func (k klogrus) With(fields ...log.Field) log.Logger {
	lfields := logrus.Fields{}

	for _, field := range fields {
		converter[field.Type](lfields, field)
	}

	new := k.core.WithFields(lfields)
	return &klogrus{new, k.name}
}

func (k klogrus) Named(s string) log.Logger {
	name := s
	if k.name != "" {
		name = strings.Join([]string{k.name, s}, ".")
	}

	new := k.core.WithField("logger", name)
	return &klogrus{new, name}
}

func (k klogrus) Debug(message string)                   { k.core.Debug(message) }
func (k klogrus) Debugf(format string, a ...interface{}) { k.Debug(fmt.Sprintf(format, a...)) }

func (k klogrus) Info(message string)                   { k.core.Info(message) }
func (k klogrus) Infof(format string, a ...interface{}) { k.Info(fmt.Sprintf(format, a...)) }

func (k klogrus) Warn(message string)                   { k.core.Warn(message) }
func (k klogrus) Warnf(format string, a ...interface{}) { k.Warn(fmt.Sprintf(format, a...)) }

func (k klogrus) Error(message string)                   { k.core.Error(message) }
func (k klogrus) Errorf(format string, a ...interface{}) { k.Error(fmt.Sprintf(format, a...)) }

func (k klogrus) Fatal(message string)                   { k.core.Fatal(message) }
func (k klogrus) Fatalf(format string, a ...interface{}) { k.Fatal(fmt.Sprintf(format, a...)) }

func (k klogrus) Panic(message string)                   { k.core.Panic(message) }
func (k klogrus) Panicf(format string, a ...interface{}) { k.Panic(fmt.Sprintf(format, a...)) }
