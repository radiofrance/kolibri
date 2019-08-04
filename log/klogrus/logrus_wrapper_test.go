package klogrus

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"

	"github.com/radiofrance/kolibri/log"
)

type klogrusTestSuite struct {
	suite.Suite
	logrus struct {
		*logrus.Logger
		out bytes.Buffer
	}
	klogrus struct {
		log.Logger
		out bytes.Buffer
	}
}

func TestNew(t *testing.T) {
	logger := logrus.New()
	klog := New(logger)

	assert.Equal(t, logger, klog.(*klogrus).core)
	assert.Equal(t, "", klog.(*klogrus).name)
}

func TestKlogrus(t *testing.T) {
	suite.Run(t, new(klogrusTestSuite))
}

func (kts *klogrusTestSuite) SetupTest() {
	kts.logrus.out = bytes.Buffer{}
	kts.logrus.Logger = logrus.New()
	kts.logrus.Logger.SetOutput(&kts.logrus.out)
	kts.logrus.Logger.SetLevel(logrus.DebugLevel)
	kts.logrus.Logger.ExitFunc = func(int) {}

	kts.klogrus.out = bytes.Buffer{}
	core := logrus.New()
	core.SetOutput(&kts.klogrus.out)
	core.SetLevel(logrus.DebugLevel)
	core.ExitFunc = func(int) {}
	kts.klogrus.Logger = New(core)
}

func (kts *klogrusTestSuite) TestWith() {
	llog, klog := kts.logrus, kts.klogrus

	addr := net.ParseIP("1.2.3.4")
	name := username("phil")
	ints := []int{5, 6}

	tests := []struct {
		name   string
		actual log.Logger
		expect logrus.FieldLogger
	}{
		{"With:Array", klog.With(log.Strings("k", []string{"ab12"})), llog.WithField("k", []string{"ab12"})},
		{"With:Binary", klog.With(log.Binary("k", []byte("ab12"))), llog.WithField("k", []byte("ab12"))},
		{"With:Bool", klog.With(log.Bool("k", true)), llog.WithField("k", true)},
		{"With:ByteString", klog.With(log.ByteString("k", []byte("ab12"))), llog.WithField("k", []byte("ab12"))},
		{"With:Complex128", klog.With(log.Complex128("k", 1+2i)), llog.WithField("k", 1+2i)},
		{"With:Complex64", klog.With(log.Complex64("k", 1+2i)), llog.WithField("k", 1+2i)},
		{"With:Duration", klog.With(log.Duration("k", 1)), llog.WithField("k", time.Duration(1))},
		{"With:Error", klog.With(log.Error("k", fmt.Errorf("err"))), llog.WithField("k", fmt.Errorf("err"))},
		{"With:Float64", klog.With(log.Float64("k", 3.14)), llog.WithField("k", 3.14)},
		{"With:Float32", klog.With(log.Float32("k", 3.14)), llog.WithField("k", 3.14)},
		{"With:Int64", klog.With(log.Int64("k", 1)), llog.WithField("k", 1)},
		{"With:Int32", klog.With(log.Int32("k", 1)), llog.WithField("k", 1)},
		{"With:Int16", klog.With(log.Int16("k", 1)), llog.WithField("k", 1)},
		{"With:Int8", klog.With(log.Int8("k", 1)), llog.WithField("k", 1)},
		{"With:Object", klog.With(log.Object("k", name)), llog.WithField("k", name)},
		{"With:Reflect", klog.With(log.Reflect("k", ints)), llog.WithField("k", ints)},
		{"With:String", klog.With(log.String("k", "foo")), llog.WithField("k", "foo")},
		{"With:Stringer", klog.With(log.Stringer("k", addr)), llog.WithField("k", addr)},
		{"With:Uint64", klog.With(log.Uint64("k", 1)), llog.WithField("k", 1)},
		{"With:Uint32", klog.With(log.Uint32("k", 1)), llog.WithField("k", 1)},
		{"With:Uint16", klog.With(log.Uint16("k", 1)), llog.WithField("k", 1)},
		{"With:Uint8", klog.With(log.Uint8("k", 1)), llog.WithField("k", 1)},
		{"With:Uintptr", klog.With(log.Uintptr("k", 0xa)), llog.WithField("k", 0xa)},
		{"With:Time", klog.With(log.Time("k", time.Unix(0, 1000).In(time.UTC))), llog.WithField("k", time.Unix(0, 1000).In(time.UTC))},
	}

	for _, tt := range tests {
		llog.out.Reset()
		klog.out.Reset()

		tt.actual.Info("_")
		tt.expect.Infof("_")
		assert.Equal(kts.T(), llog.out.String(), klog.out.String(), "Unexpected output from %s.", tt.name)
	}

}

func (kts *klogrusTestSuite) TestNamed() {
	llog, klog := kts.logrus, kts.klogrus

	tests := []struct {
		name   string
		actual log.Logger
		expect logrus.FieldLogger
	}{
		{"Named:a", klog.Named("a"), llog.WithField("logger", "a")},
		{"Named:a.b", klog.Named("a").Named("b"), llog.WithField("logger", "a.b")},
		{"Named:a.b.c", klog.Named("a").Named("b").Named("c"), llog.WithField("logger", "a.b.c")},
		{"Named:a.b.c.d", klog.Named("a").Named("b").Named("c").Named("d"), llog.WithField("logger", "a.b.c.d")},
	}

	for _, tt := range tests {
		llog.out.Reset()
		klog.out.Reset()

		tt.actual.Info("_")
		tt.expect.Infof("_")
		assert.Equal(kts.T(), llog.out.String(), klog.out.String(), "Unexpected output from %s.", tt.name)
	}
}

func (kts *klogrusTestSuite) TestLog() {
	llog, klog := kts.logrus, kts.klogrus

	tests := []struct {
		name   string
		actual func(string)
		expect func(...interface{})
	}{
		{"Level:Debug", klog.Debug, llog.Debug},
		{"Level:Info", klog.Info, llog.Info},
		{"Level:Warn", klog.Warn, llog.Warn},
		{"Level:Error", klog.Error, llog.Error},
		{"Level:Fatal", klog.Fatal, llog.Fatal},
	}

	for _, tt := range tests {
		llog.out.Reset()
		klog.out.Reset()

		tt.actual("_")
		tt.expect("_")
		assert.Equal(kts.T(), llog.out.String(), klog.out.String(), "Unexpected output from %s.", tt.name)
	}
}

func (kts *klogrusTestSuite) TestLogFormat() {
	llog, klog := kts.logrus, kts.klogrus

	tests := []struct {
		name   string
		actual func(string, ...interface{})
		expect func(string, ...interface{})
	}{
		{"Level:Debug", klog.Debugf, llog.Debugf},
		{"Level:Info", klog.Infof, llog.Infof},
		{"Level:Warn", klog.Warnf, llog.Warnf},
		{"Level:Error", klog.Errorf, llog.Errorf},
		{"Level:Fatal", klog.Fatalf, llog.Fatalf},
	}

	for _, tt := range tests {
		llog.out.Reset()
		klog.out.Reset()

		tt.actual("_%s", "_")
		tt.expect("_%s", "_")
		assert.Equal(kts.T(), llog.out.String(), klog.out.String(), "Unexpected output from %s.", tt.name)
	}
}

func (kts *klogrusTestSuite) TestPanic() {
	llog, klog := kts.logrus, kts.klogrus

	ignorePanic := func() { recover() }
	mustPanic := func(s string) {
		assert.NotNil(kts.T(), recover(), "klogrus.%s() must panic.", s)
	}

	llog.out.Reset()
	klog.out.Reset()
	func() { defer ignorePanic(); llog.Panic("_") }()
	func() { defer mustPanic("Panic"); klog.Panic("_") }()
	assert.Equal(kts.T(), llog.out.String(), klog.out.String(), "Unexpected output from Level:Panic.")

	llog.out.Reset()
	klog.out.Reset()
	func() { defer ignorePanic(); llog.Panicf("_%s", "_") }()
	func() { defer mustPanic("Panicf"); klog.Panicf("_%s", "_") }()
	assert.Equal(kts.T(), llog.out.String(), klog.out.String(), "Unexpected output from Level:Panic.")
}
func (kts *klogrusTestSuite) TestAllInOne() {
	llog, klog := kts.logrus, kts.klogrus

	klogA := klog.Named("A")
	klogB := klog.Named("B")
	klogBC := klogB.Named("C").With(log.Duration("wait", time.Second))

	llog.out.Reset()
	klog.out.Reset()
	klogA.
		With(log.Any("err", fmt.Errorf("unknownErr"))).
		Error("Something went wrong")
	llog.
		WithField("logger", "A").
		WithField("err", fmt.Errorf("unknownErr")).
		Error("Something went wrong")
	assert.Equal(kts.T(), llog.out.String(), klog.out.String())

	llog.out.Reset()
	klog.out.Reset()
	klogB.
		Debug("...")
	llog.
		WithField("logger", "B").
		Debug("...")
	assert.Equal(kts.T(), llog.out.String(), klog.out.String())

	llog.out.Reset()
	klog.out.Reset()
	klogBC.
		With(log.Strings("___", []string{"aaa", "bbb"})).
		Warn("")
	llog.
		WithField("logger", "B.C").
		WithField("wait", time.Second).
		WithField("___", []string{"aaa", "bbb"}).
		Warn("")
	assert.Equal(kts.T(), llog.out.String(), klog.out.String())
}

type username string

func (n username) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("username", string(n))
	return nil
}
