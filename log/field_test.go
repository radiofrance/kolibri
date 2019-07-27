package log

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type username string

func (n username) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("username", string(n))
	return nil
}

func TestFieldConstructors(t *testing.T) {
	addr := net.ParseIP("1.2.3.4")
	name := username("phil")
	ints := []int{5, 6}

	tests := []struct {
		name   string
		actual Field
		expect zap.Field
	}{
		{"Binary", Binary("k", []byte("ab12")), zap.Binary("k", []byte("ab12"))},
		{"Bool", Bool("k", true), zap.Bool("k", true)},
		{"ByteString", ByteString("k", []byte("ab12")), zap.ByteString("k", []byte("ab12"))},
		{"Complex128", Complex128("k", 1+2i), zap.Complex128("k", 1+2i)},
		{"Complex64", Complex64("k", 1+2i), zap.Complex64("k", 1+2i)},
		{"Duration", Duration("k", 1), zap.Duration("k", 1)},
		{"Error", Error("k", fmt.Errorf("err")), zap.NamedError("k", fmt.Errorf("err"))},
		{"Float64", Float64("k", 3.14), zap.Float64("k", 3.14)},
		{"Float32", Float32("k", 3.14), zap.Float32("k", 3.14)},
		{"Int", Int("k", 1), zap.Int("k", 1)},
		{"Int64", Int64("k", 1), zap.Int64("k", 1)},
		{"Int32", Int32("k", 1), zap.Int32("k", 1)},
		{"Int16", Int16("k", 1), zap.Int16("k", 1)},
		{"Int8", Int8("k", 1), zap.Int8("k", 1)},
		{"Object", Object("k", name), zap.Object("k", name)},
		{"Reflect", Reflect("k", ints), zap.Reflect("k", ints)},
		{"String", String("k", "foo"), zap.String("k", "foo")},
		{"Stringer", Stringer("k", addr), zap.Stringer("k", addr)},
		{"Uint", Uint("k", 1), zap.Uint("k", 1)},
		{"Uint64", Uint64("k", 1), zap.Uint64("k", 1)},
		{"Uint32", Uint32("k", 1), zap.Uint32("k", 1)},
		{"Uint16", Uint16("k", 1), zap.Uint16("k", 1)},
		{"Uint8", Uint8("k", 1), zap.Uint8("k", 1)},
		{"Uintptr", Uintptr("k", 0xa), zap.Uintptr("k", 0xa)},
		{"Time", Time("k", time.Unix(0, 0).In(time.UTC)), zap.Time("k", time.Unix(0, 0).In(time.UTC))},
		{"Time", Time("k", time.Unix(0, 1000).In(time.UTC)), zap.Time("k", time.Unix(0, 1000).In(time.UTC))},

		{"Bools", Bools("k", []bool{true}), zap.Bools("k", []bool{true})},
		{"ByteStrings", ByteStrings("k", [][]byte{[]byte("ab12")}), zap.ByteStrings("k", [][]byte{[]byte("ab12")})},
		{"Complex128s", Complex128s("k", []complex128{1 + 2i}), zap.Complex128s("k", []complex128{1 + 2i})},
		{"Complex64s", Complex64s("k", []complex64{1 + 2i}), zap.Complex64s("k", []complex64{1 + 2i})},
		{"Durations", Durations("k", []time.Duration{1}), zap.Durations("k", []time.Duration{1})},
		{"Errors", Errors("k", []error{fmt.Errorf("err")}), zap.Errors("k", []error{fmt.Errorf("err")})},
		{"Float64s", Float64s("k", []float64{3.14}), zap.Float64s("k", []float64{3.14})},
		{"Float32s", Float32s("k", []float32{3.14}), zap.Float32s("k", []float32{3.14})},
		{"Ints", Ints("k", []int{1}), zap.Ints("k", []int{1})},
		{"Int64s", Int64s("k", []int64{1}), zap.Int64s("k", []int64{1})},
		{"Int32s", Int32s("k", []int32{1}), zap.Int32s("k", []int32{1})},
		{"Int16s", Int16s("k", []int16{1}), zap.Int16s("k", []int16{1})},
		{"Int8s", Int8s("k", []int8{1}), zap.Int8s("k", []int8{1})},
		{"Strings", Strings("k", []string{"foo"}), zap.Strings("k", []string{"foo"})},
		{"Times", Times("k", []time.Time{time.Unix(0, 1000).In(time.UTC)}), zap.Times("k", []time.Time{time.Unix(0, 1000).In(time.UTC)})},
		{"Uints", Uints("k", []uint{1}), zap.Uints("k", []uint{1})},
		{"Uint64s", Uint64s("k", []uint64{1}), zap.Uint64s("k", []uint64{1})},
		{"Uint32s", Uint32s("k", []uint32{1}), zap.Uint32s("k", []uint32{1})},
		{"Uint16s", Uint16s("k", []uint16{1}), zap.Uint16s("k", []uint16{1})},
		{"Uint8s", Uint8s("k", []uint8{1}), zap.Uint8s("k", []uint8{1})},
		{"Uintptrs", Uintptrs("k", []uintptr{0xa}), zap.Uintptrs("k", []uintptr{0xa})},

		{"Any:Binary", Any("k", []byte("ab12")), zap.Any("k", []byte("ab12"))},
		{"Any:Bool", Any("k", true), zap.Any("k", true)},
		{"Any:ByteString", Any("k", []byte("ab12")), zap.Any("k", []byte("ab12"))},
		{"Any:Complex128", Any("k", 1+2i), zap.Any("k", 1+2i)},
		{"Any:Complex64", Any("k", 1+2i), zap.Any("k", 1+2i)},
		{"Any:Duration", Any("k", 1), zap.Any("k", 1)},
		{"Any:Error", Any("k", fmt.Errorf("err")), zap.Any("k", fmt.Errorf("err"))},
		{"Any:Float64", Any("k", 3.14), zap.Any("k", 3.14)},
		{"Any:Float32", Any("k", 3.14), zap.Any("k", 3.14)},
		{"Any:Int", Any("k", 1), zap.Any("k", 1)},
		{"Any:Int64", Any("k", 1), zap.Any("k", 1)},
		{"Any:Int32", Any("k", 1), zap.Any("k", 1)},
		{"Any:Int16", Any("k", 1), zap.Any("k", 1)},
		{"Any:Int8", Any("k", 1), zap.Any("k", 1)},
		{"Any:Object", Any("k", name), zap.Any("k", name)},
		{"Any:Reflect", Any("k", ints), zap.Any("k", ints)},
		{"Any:String", Any("k", "foo"), zap.Any("k", "foo")},
		{"Any:Stringer", Any("k", addr), zap.Any("k", addr)},
		{"Any:Uint", Any("k", 1), zap.Any("k", 1)},
		{"Any:Uint64", Any("k", 1), zap.Any("k", 1)},
		{"Any:Uint32", Any("k", 1), zap.Any("k", 1)},
		{"Any:Uint16", Any("k", 1), zap.Any("k", 1)},
		{"Any:Uint8", Any("k", 1), zap.Any("k", 1)},
		{"Any:Uintptr", Any("k", 0xa), zap.Any("k", 0xa)},
		{"Any:Time", Any("k", time.Unix(0, 0).In(time.UTC)), zap.Any("k", time.Unix(0, 0).In(time.UTC))},
		{"Any:Time", Any("k", time.Unix(0, 1000).In(time.UTC)), zap.Any("k", time.Unix(0, 1000).In(time.UTC))},

		{"Any:Bools", Any("k", []bool{true}), zap.Any("k", []bool{true})},
		{"Any:ByteStrings", Any("k", [][]byte{[]byte("ab12")}), zap.Any("k", [][]byte{[]byte("ab12")})},
		{"Any:Complex128s", Any("k", []complex128{1 + 2i}), zap.Any("k", []complex128{1 + 2i})},
		{"Any:Complex64s", Any("k", []complex64{1 + 2i}), zap.Any("k", []complex64{1 + 2i})},
		{"Any:Durations", Any("k", []time.Duration{1}), zap.Any("k", []time.Duration{1})},
		{"Any:Errors", Any("k", []error{fmt.Errorf("err")}), zap.Any("k", []error{fmt.Errorf("err")})},
		{"Any:Float64s", Any("k", []float64{3.14}), zap.Any("k", []float64{3.14})},
		{"Any:Float32s", Any("k", []float32{3.14}), zap.Any("k", []float32{3.14})},
		{"Any:Ints", Any("k", []int{1}), zap.Any("k", []int{1})},
		{"Any:Int64s", Any("k", []int64{1}), zap.Any("k", []int64{1})},
		{"Any:Int32s", Any("k", []int32{1}), zap.Any("k", []int32{1})},
		{"Any:Int16s", Any("k", []int16{1}), zap.Any("k", []int16{1})},
		{"Any:Int8s", Any("k", []int8{1}), zap.Any("k", []int8{1})},
		{"Any:Strings", Any("k", []string{"foo"}), zap.Any("k", []string{"foo"})},
		{"Any:Times", Any("k", []time.Time{time.Unix(0, 1000).In(time.UTC)}), zap.Any("k", []time.Time{time.Unix(0, 1000).In(time.UTC)})},
		{"Any:Uints", Any("k", []uint{1}), zap.Any("k", []uint{1})},
		{"Any:Uint64s", Any("k", []uint64{1}), zap.Any("k", []uint64{1})},
		{"Any:Uint32s", Any("k", []uint32{1}), zap.Any("k", []uint32{1})},
		{"Any:Uint16s", Any("k", []uint16{1}), zap.Any("k", []uint16{1})},
		{"Any:Uint8s", Any("k", []uint8{1}), zap.Any("k", []uint8{1})},
		{"Any:Uintptrs", Any("k", []uintptr{0xa}), zap.Any("k", []uintptr{0xa})},
	}

	for _, tt := range tests {
		if !assert.Equal(t, tt.actual, Field(tt.expect), "Unexpected output from convenience field constructor %s.", tt.name) {
			t.Logf("type expected: %T\nGot: %T", tt.expect.Interface, tt.actual.Interface)
		}
	}
}

func TestStackField(t *testing.T) {
	f, zf := Stack("stacktrace"), zap.Stack("stacktrace")

	assert.Equal(t, "stacktrace", f.Key, "Unexpected field key.")
	assert.Equal(t, zapcore.StringType, f.Type, "Unexpected field type.")
	assert.Equal(t, zf.String, f.String, "Unexpected stack trace.") // Tested by zap.Stack test
}
