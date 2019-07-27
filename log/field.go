package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, val []byte) Field { return Field(zap.Binary(key, val)) }

// Bool constructs a field that carries a bool.
func Bool(key string, val bool) Field { return Field(zap.Bool(key, val)) }

// ByteString constructs a field that carries UTF-8 encoded text as a []byte.
// To log opaque binary blobs (which aren't necessarily valid UTF-8), use
// Binary.
func ByteString(key string, val []byte) Field { return Field(zap.ByteString(key, val)) }

// Complex128 constructs a field that carries a complex number. Unlike most
// numeric fields, this costs an allocation (to convert the complex128 to
// interface{}).
func Complex128(key string, val complex128) Field { return Field(zap.Complex128(key, val)) }

// Complex64 constructs a field that carries a complex number. Unlike most
// numeric fields, this costs an allocation (to convert the complex64 to
// interface{}).
func Complex64(key string, val complex64) Field { return Field(zap.Complex64(key, val)) }

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func Duration(key string, val time.Duration) Field { return Field(zap.Duration(key, val)) }

// Error constructs a field that lazily stores err.Error() under the
// provided key. Errors which also implement fmt.Formatter (like those produced
// by github.com/pkg/errors) will also have their verbose representation stored
// under key+"Verbose". If passed a nil error, the field is a no-op.
//
// For the common case in which the key is simply "error", the Error function
// is shorter and less repetitive.
func Error(key string, val error) Field { return Field(zap.NamedError(key, val)) }

// Float64 constructs a field that carries a float64. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float64(key string, val float64) Field { return Field(zap.Float64(key, val)) }

// Float32 constructs a field that carries a float32. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float32(key string, val float32) Field { return Field(zap.Float32(key, val)) }

// Int constructs a field with the given key and value.
func Int(key string, val int) Field { return Field(zap.Int(key, val)) }

// Int64 constructs a field with the given key and value.
func Int64(key string, val int64) Field { return Field(zap.Int64(key, val)) }

// Int32 constructs a field with the given key and value.
func Int32(key string, val int32) Field { return Field(zap.Int32(key, val)) }

// Int16 constructs a field with the given key and value.
func Int16(key string, val int16) Field { return Field(zap.Int16(key, val)) }

// Int8 constructs a field with the given key and value.
func Int8(key string, val int8) Field { return Field(zap.Int8(key, val)) }

// Object constructs a field with the given key and ObjectMarshaler. It
// provides a flexible, but still type-safe and efficient, way to add map- or
// struct-like user-defined types to the logging context.
func Object(key string, val zapcore.ObjectMarshaler) Field { return Field(zap.Object(key, val)) }

// Reflect constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy. Outside tests, Any is always a better choice.
//
// If encoding fails (e.g., trying to serialize a map[int]string to JSON), Reflect
// includes the error message in the final log output.
func Reflect(key string, val interface{}) Field { return Field(zap.Reflect(key, val)) }

// Stack constructs a field that stores a stacktrace of the current goroutine
// under provided key. Keep in mind that taking a stacktrace is eager and
// expensive (relatively speaking); this function both makes an allocation and
// takes about two microseconds.
func Stack(key string) Field {
	field := Field(zap.Stack(key))
	stack := field.String

	// Remove current function in stacktrace
	n := 2
	for i, c := range stack {
		if c == '\n' {
			n--
		}
		if n == 0 {
			field.String = stack[i+1:]
			break
		}
	}

	return field
}

// String constructs a field with the given key and value.
func String(key string, val string) Field { return Field(zap.String(key, val)) }

// Stringer constructs a field with the given key and the output of the value's
// String method. The Stringer's String method is called lazily.
func Stringer(key string, val fmt.Stringer) Field { return Field(zap.Stringer(key, val)) }

// Time constructs a Field with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) Field { return Field(zap.Time(key, val)) }

// Uint constructs a field with the given key and value.
func Uint(key string, val uint) Field { return Field(zap.Uint(key, val)) }

// Uint64 constructs a field with the given key and value.
func Uint64(key string, val uint64) Field { return Field(zap.Uint64(key, val)) }

// Uint32 constructs a field with the given key and value.
func Uint32(key string, val uint32) Field { return Field(zap.Uint32(key, val)) }

// Uint16 constructs a field with the given key and value.
func Uint16(key string, val uint16) Field { return Field(zap.Uint16(key, val)) }

// Uint8 constructs a field with the given key and value.
func Uint8(key string, val uint8) Field { return Field(zap.Uint8(key, val)) }

// Uintptr constructs a field with the given key and value.
func Uintptr(key string, val uintptr) Field { return Field(zap.Uintptr(key, val)) }

// Bools constructs a field that carries a slice of bools.
func Bools(key string, val []bool) Field { return Field(zap.Bools(key, val)) }

// ByteStrings constructs a field that carries a slice of []byte, each of which
// must be UTF-8 encoded text.
func ByteStrings(key string, val [][]byte) Field { return Field(zap.ByteStrings(key, val)) }

// Complex128s constructs a field that carries a slice of complex numbers.
func Complex128s(key string, val []complex128) Field { return Field(zap.Complex128s(key, val)) }

// Complex64s constructs a field that carries a slice of complex numbers.
func Complex64s(key string, val []complex64) Field { return Field(zap.Complex64s(key, val)) }

// Durations constructs a field that carries a slice of time.Durations.
func Durations(key string, val []time.Duration) Field { return Field(zap.Durations(key, val)) }

// Errors constructs a field that carries a slice of errors.
func Errors(key string, val []error) Field { return Field(zap.Errors(key, val)) }

// Float64s constructs a field that carries a slice of floats.
func Float64s(key string, val []float64) Field { return Field(zap.Float64s(key, val)) }

// Float32s constructs a field that carries a slice of floats.
func Float32s(key string, val []float32) Field { return Field(zap.Float32s(key, val)) }

// Ints constructs a field that carries a slice of integers.
func Ints(key string, val []int) Field { return Field(zap.Ints(key, val)) }

// Int64s constructs a field that carries a slice of integers.
func Int64s(key string, val []int64) Field { return Field(zap.Int64s(key, val)) }

// Int32s constructs a field that carries a slice of integers.
func Int32s(key string, val []int32) Field { return Field(zap.Int32s(key, val)) }

// Int16s constructs a field that carries a slice of integers.
func Int16s(key string, val []int16) Field { return Field(zap.Int16s(key, val)) }

// Int8s constructs a field that carries a slice of integers.
func Int8s(key string, val []int8) Field { return Field(zap.Int8s(key, val)) }

// Strings constructs a field that carries a slice of strings.
func Strings(key string, val []string) Field { return Field(zap.Strings(key, val)) }

// Times constructs a field that carries a slice of time.Times.
func Times(key string, val []time.Time) Field { return Field(zap.Times(key, val)) }

// Uints constructs a field that carries a slice of unsigned integers.
func Uints(key string, val []uint) Field { return Field(zap.Uints(key, val)) }

// Uint64s constructs a field that carries a slice of unsigned integers.
func Uint64s(key string, val []uint64) Field { return Field(zap.Uint64s(key, val)) }

// Uint32s constructs a field that carries a slice of unsigned integers.
func Uint32s(key string, val []uint32) Field { return Field(zap.Uint32s(key, val)) }

// Uint16s constructs a field that carries a slice of unsigned integers.
func Uint16s(key string, val []uint16) Field { return Field(zap.Uint16s(key, val)) }

// Uint8s constructs a field that carries a slice of unsigned integers.
func Uint8s(key string, val []uint8) Field { return Field(zap.Uint8s(key, val)) }

// Uintptrs constructs a field that carries a slice of pointer addresses.
func Uintptrs(key string, val []uintptr) Field { return Field(zap.Uintptrs(key, val)) }

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value interface{}) Field { return Field(zap.Any(key, value)) }
