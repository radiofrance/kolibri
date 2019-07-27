package kzap

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gitlab.com/xunleii.io/kolibri/log"
)

type kzapTestSuite struct {
	suite.Suite
	zap struct {
		*zap.Logger
		out wsBuffer
	}
	kzap struct {
		log.Logger
		out wsBuffer
	}
}

func TestNew(t *testing.T) {
	logger, _ := zap.NewProduction()
	klog := New(logger)

	assert.Equal(t, logger, klog.(*kzap).core)
}

func TestKzap(t *testing.T) {
	suite.Run(t, new(kzapTestSuite))
}

func (kts *kzapTestSuite) SetupTest() {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = func(_ time.Time, enc zapcore.PrimitiveArrayEncoder) { enc.AppendFloat64(0) }
	enc := zapcore.NewJSONEncoder(cfg)
	enab := zap.LevelEnablerFunc(func(l zapcore.Level) bool { return true })

	kts.zap.out = wsBuffer{&bytes.Buffer{}}
	kts.zap.Logger = zap.New(zapcore.NewCore(enc, kts.zap.out, enab))
	kts.kzap.out = wsBuffer{&bytes.Buffer{}}
	kts.kzap.Logger = New(zap.New(zapcore.NewCore(enc, kts.kzap.out, enab)))
}

func (kts *kzapTestSuite) TestAllInOne() {
	zlog, klog := kts.zap, kts.kzap

	klogA := klog.Named("A")
	klogB := klog.Named("B")
	klogBC := klogB.Named("C").With(log.Duration("wait", time.Second))

	zlog.out.Reset()
	klog.out.Reset()
	klogA.
		With(log.Any("err", fmt.Errorf("unknownErr"))).
		Error("Something went wrong")
	zlog.
		Named("A").
		With(zap.Any("err", fmt.Errorf("unknownErr"))).
		Error("Something went wrong")
	assert.Equal(kts.T(), zlog.out.String(), klog.out.String())

	zlog.out.Reset()
	klog.out.Reset()
	klogB.
		Debug("...")
	zlog.
		Named("B").
		Debug("...")
	assert.Equal(kts.T(), zlog.out.String(), klog.out.String())

	zlog.out.Reset()
	klog.out.Reset()
	klogBC.
		With(log.Strings("___", []string{"aaa", "bbb"})).
		Warn("")
	zlog.
		Named("B").
		Named("C").
		With(zap.Duration("wait", time.Second)).
		With(zap.Strings("___", []string{"aaa", "bbb"})).
		Warn("")
	assert.Equal(kts.T(), zlog.out.String(), klog.out.String())
}

func (kts *kzapTestSuite) TestKzap_With() {
	zlog, klog := kts.zap, kts.kzap

	addr := net.ParseIP("1.2.3.4")
	name := username("phil")
	ints := []int{5, 6}

	tests := []struct {
		name   string
		actual log.Logger
		expect *zap.Logger
	}{
		{"With:Array", klog.With(log.Strings("k", []string{"ab12"})), zlog.With(zap.Strings("k", []string{"ab12"}))},
		{"With:Binary", klog.With(log.Binary("k", []byte("ab12"))), zlog.With(zap.Binary("k", []byte("ab12")))},
		{"With:Bool", klog.With(log.Bool("k", true)), zlog.With(zap.Bool("k", true))},
		{"With:ByteString", klog.With(log.ByteString("k", []byte("ab12"))), zlog.With(zap.ByteString("k", []byte("ab12")))},
		{"With:Complex128", klog.With(log.Complex128("k", 1+2i)), zlog.With(zap.Complex128("k", 1+2i))},
		{"With:Complex64", klog.With(log.Complex64("k", 1+2i)), zlog.With(zap.Complex64("k", 1+2i))},
		{"With:Duration", klog.With(log.Duration("k", 1)), zlog.With(zap.Duration("k", 1))},
		{"With:Error", klog.With(log.Error("k", fmt.Errorf("err"))), zlog.With(zap.NamedError("k", fmt.Errorf("err")))},
		{"With:Float64", klog.With(log.Float64("k", 3.14)), zlog.With(zap.Float64("k", 3.14))},
		{"With:Float32", klog.With(log.Float32("k", 3.14)), zlog.With(zap.Float32("k", 3.14))},
		{"With:Int64", klog.With(log.Int64("k", 1)), zlog.With(zap.Int64("k", 1))},
		{"With:Int32", klog.With(log.Int32("k", 1)), zlog.With(zap.Int32("k", 1))},
		{"With:Int16", klog.With(log.Int16("k", 1)), zlog.With(zap.Int16("k", 1))},
		{"With:Int8", klog.With(log.Int8("k", 1)), zlog.With(zap.Int8("k", 1))},
		{"With:Object", klog.With(log.Object("k", name)), zlog.With(zap.Object("k", name))},
		{"With:Reflect", klog.With(log.Reflect("k", ints)), zlog.With(zap.Reflect("k", ints))},
		{"With:String", klog.With(log.String("k", "foo")), zlog.With(zap.String("k", "foo"))},
		{"With:Stringer", klog.With(log.Stringer("k", addr)), zlog.With(zap.Stringer("k", addr))},
		{"With:Uint64", klog.With(log.Uint64("k", 1)), zlog.With(zap.Uint64("k", 1))},
		{"With:Uint32", klog.With(log.Uint32("k", 1)), zlog.With(zap.Uint32("k", 1))},
		{"With:Uint16", klog.With(log.Uint16("k", 1)), zlog.With(zap.Uint16("k", 1))},
		{"With:Uint8", klog.With(log.Uint8("k", 1)), zlog.With(zap.Uint8("k", 1))},
		{"With:Uintptr", klog.With(log.Uintptr("k", 0xa)), zlog.With(zap.Uintptr("k", 0xa))},
		{"With:Time", klog.With(log.Time("k", time.Unix(0, 1000).In(time.UTC))), zlog.With(zap.Time("k", time.Unix(0, 1000).In(time.UTC)))},
	}

	for _, tt := range tests {
		zlog.out.Reset()
		klog.out.Reset()

		tt.actual.Info("_")
		tt.expect.Info("_")
		assert.Equal(kts.T(), zlog.out.String(), klog.out.String(), "Unexpected output from %s.", tt.name)
	}

}

func (kts *kzapTestSuite) TestKlogrus_Named() {
	zlog, klog := kts.zap, kts.kzap

	tests := []struct {
		name   string
		actual log.Logger
		expect *zap.Logger
	}{
		{"Named:a", klog.Named("a"), zlog.Named("a")},
		{"Named:a.b", klog.Named("a").Named("b"), zlog.Named("a").Named("b")},
		{"Named:a.b.c", klog.Named("a").Named("b").Named("c"), zlog.Named("a").Named("b").Named("c")},
		{"Named:a.b.c.d", klog.Named("a").Named("b").Named("c").Named("d"), zlog.Named("a").Named("b").Named("c").Named("d")},
	}

	for _, tt := range tests {
		zlog.out.Reset()
		klog.out.Reset()

		tt.actual.Info("_")
		tt.expect.Info("_")
		assert.Equal(kts.T(), zlog.out.String(), klog.out.String(), "Unexpected output from %s.", tt.name)
	}
}

func (kts *kzapTestSuite) TestKlogrus_Log() {
	zlog, klog := kts.zap, kts.kzap

	tests := []struct {
		name   string
		actual func(string)
		expect func(string, ...zap.Field)
	}{
		{"Level:Debug", klog.Debug, zlog.Debug},
		{"Level:Info", klog.Info, zlog.Info},
		{"Level:Warn", klog.Warn, zlog.Warn},
		{"Level:Error", klog.Error, zlog.Error},
	}

	for _, tt := range tests {
		zlog.out.Reset()
		klog.out.Reset()

		tt.actual("_")
		tt.expect("_")
		assert.Equal(kts.T(), zlog.out.String(), klog.out.String(), "Unexpected output from %s.", tt.name)
	}

}

func (kts *kzapTestSuite) TestKlogrus_LogFormat() {
	zlog, klog := kts.zap, kts.kzap

	tests := []struct {
		name   string
		actual func(string, ...interface{})
		expect func(string, ...zap.Field)
	}{
		{"Level:Debug", klog.Debugf, zlog.Debug},
		{"Level:Info", klog.Infof, zlog.Info},
		{"Level:Warn", klog.Warnf, zlog.Warn},
		{"Level:Error", klog.Errorf, zlog.Error},
	}

	for _, tt := range tests {
		zlog.out.Reset()
		klog.out.Reset()

		tt.actual("_%s", "_")
		tt.expect(fmt.Sprintf("_%s", "_"))
		assert.Equal(kts.T(), zlog.out.String(), klog.out.String(), "Unexpected output from %s.", tt.name)
	}
}

func (kts *kzapTestSuite) TestKlogrus_Panic() {
	zlog, klog := kts.zap, kts.kzap

	ignorePanic := func() { recover() }
	mustPanic := func(s string) {
		assert.NotNil(kts.T(), recover(), "kzap.%s() must panic.", s)
	}

	zlog.out.Reset()
	klog.out.Reset()
	func() { defer ignorePanic(); zlog.Panic("_") }()
	func() { defer mustPanic("Panic"); klog.Panic("_") }()
	assert.Equal(kts.T(), zlog.out.String(), klog.out.String(), "Unexpected output from Level:Panic.")

	zlog.out.Reset()
	klog.out.Reset()
	func() { defer ignorePanic(); zlog.Panic(fmt.Sprintf("_%s", "_")) }()
	func() { defer mustPanic("Panicf"); klog.Panicf("_%s", "_") }()
	assert.Equal(kts.T(), zlog.out.String(), klog.out.String(), "Unexpected output from Level:Panic.")
}

type username string

func (n username) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("username", string(n))
	return nil
}

type wsBuffer struct{ *bytes.Buffer }

func (wsBuffer) Sync() error { return nil }
