package flog

import (
	"bytes"
	"math"
	"testing"
	"unsafe"

	"github.com/aidenfine/flog/src"
)

func fieldConsole(t *testing.T, field Field) string {
	t.Helper()
	var out bytes.Buffer
	if _, err := new(src.Writer).WriteLog(&out, &src.Entity{Fields: []Field{field}}, "", "console"); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func TestFieldHelpers(t *testing.T) {
	number := 42
	ch := make(chan int)
	fn := func() {}
	tests := []struct {
		name  string
		field Field
		value any
	}{
		{"null", Null("key"), nil},
		{"string", String("key", "hello"), "hello"},
		{"bool true", Bool("key", true), true},
		{"bool false", Bool("key", false), false},
		{"int", Int("key", -42), int(-42)},
		{"int8", Int8("key", -128), int8(-128)},
		{"int16", Int16("key", -32768), int16(-32768)},
		{"int32", Int32("key", math.MinInt32), int32(math.MinInt32)},
		{"int64", Int64("key", math.MinInt64), int64(math.MinInt64)},
		{"uint", Uint("key", ^uint(0)), ^uint(0)},
		{"uint8", Uint8("key", math.MaxUint8), uint8(math.MaxUint8)},
		{"uint16", Uint16("key", math.MaxUint16), uint16(math.MaxUint16)},
		{"uint16 alias", Unit16("key", math.MaxUint16), uint16(math.MaxUint16)},
		{"uint32", Uint32("key", math.MaxUint32), uint32(math.MaxUint32)},
		{"uint64", Uint64("key", math.MaxUint64), uint64(math.MaxUint64)},
		{"uintptr", Uintptr("key", ^uintptr(0)), ^uintptr(0)},
		{"float32", Float32("key", 1.25), float32(1.25)},
		{"float64", Float64("key", -2.5), float64(-2.5)},
		{"complex128", Complex128("key", 1+2i), complex128(1 + 2i)},
		{"binary", Binary("key", []byte{0, 255}), []byte{0, 255}},
		{"array", Array("key", [2]int{1, 2}), [2]int{1, 2}},
		{"slice", Array("key", []int{1, 2}), []int{1, 2}},
		{"object", Object("key", struct{ Count int }{3}), struct{ Count int }{3}},
		{"map", Object("key", map[string]int{"a": 1}), map[string]int{"a": 1}},
		{"pointer", Pointer("key", &number), &number},
		{"nil pointer", Pointer[int]("key", nil), (*int)(nil)},
		{"channel", Chan("key", ch), ch},
		{"receive channel", Chan("key", (<-chan int)(ch)), (<-chan int)(ch)},
		{"send channel", Chan("key", (chan<- int)(ch)), (chan<- int)(ch)},
		{"function", Func("key", fn), fn},
		{"unsafe pointer", UnsafePointer("key", unsafe.Pointer(&number)), unsafe.Pointer(&number)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fieldConsole(t, tt.field)
			want := fieldConsole(t, Any("key", tt.value))
			if got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
	if got := fieldConsole(t, ByteString("key", []byte("hello"))); got != "{\"key\": \"hello\"}\n" {
		t.Fatal(got)
	}
}

func TestInvalidFieldHelpers(t *testing.T) {
	for _, field := range []Field{Invalid("key", 42), Array("key", 42), Object("key", 42), Chan("key", 42), Func("key", 42)} {
		var out bytes.Buffer
		if _, err := new(src.Writer).WriteLog(&out, &src.Entity{Fields: []Field{field}}, "", "console"); err == nil || out.Len() != 0 {
			t.Fatalf("field=%+v error=%v output=%q", field, err, out.String())
		}
	}
}

func TestPointerFieldHelpers(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		var value string = "hello"
		got := fieldConsole(t, StringPtr("key", &value))
		want := fieldConsole(t, String("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *string
		if got := fieldConsole(t, StringPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Bool", func(t *testing.T) {
		var value bool = true
		got := fieldConsole(t, BoolPtr("key", &value))
		want := fieldConsole(t, Bool("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *bool
		if got := fieldConsole(t, BoolPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Int", func(t *testing.T) {
		var value int = -42
		got := fieldConsole(t, IntPtr("key", &value))
		want := fieldConsole(t, Int("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *int
		if got := fieldConsole(t, IntPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Int8", func(t *testing.T) {
		var value int8 = -12
		got := fieldConsole(t, Int8Ptr("key", &value))
		want := fieldConsole(t, Int8("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *int8
		if got := fieldConsole(t, Int8Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Int16", func(t *testing.T) {
		var value int16 = -123
		got := fieldConsole(t, Int16Ptr("key", &value))
		want := fieldConsole(t, Int16("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *int16
		if got := fieldConsole(t, Int16Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Int32", func(t *testing.T) {
		var value int32 = -1234
		got := fieldConsole(t, Int32Ptr("key", &value))
		want := fieldConsole(t, Int32("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *int32
		if got := fieldConsole(t, Int32Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Int64", func(t *testing.T) {
		var value int64 = -12345
		got := fieldConsole(t, Int64Ptr("key", &value))
		want := fieldConsole(t, Int64("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *int64
		if got := fieldConsole(t, Int64Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Uint", func(t *testing.T) {
		var value uint = 42
		got := fieldConsole(t, UintPtr("key", &value))
		want := fieldConsole(t, Uint("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *uint
		if got := fieldConsole(t, UintPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Uint8", func(t *testing.T) {
		var value uint8 = 255
		got := fieldConsole(t, Uint8Ptr("key", &value))
		want := fieldConsole(t, Uint8("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *uint8
		if got := fieldConsole(t, Uint8Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Uint16", func(t *testing.T) {
		var value uint16 = 65535
		got := fieldConsole(t, Uint16Ptr("key", &value))
		want := fieldConsole(t, Uint16("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *uint16
		if got := fieldConsole(t, Uint16Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Uint32", func(t *testing.T) {
		var value uint32 = math.MaxUint32
		got := fieldConsole(t, Uint32Ptr("key", &value))
		want := fieldConsole(t, Uint32("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *uint32
		if got := fieldConsole(t, Uint32Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Uint64", func(t *testing.T) {
		var value uint64 = math.MaxUint64
		got := fieldConsole(t, Uint64Ptr("key", &value))
		want := fieldConsole(t, Uint64("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *uint64
		if got := fieldConsole(t, Uint64Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Uintptr", func(t *testing.T) {
		var value uintptr = 42
		got := fieldConsole(t, UintptrPtr("key", &value))
		want := fieldConsole(t, Uintptr("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *uintptr
		if got := fieldConsole(t, UintptrPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Float32", func(t *testing.T) {
		var value float32 = 1.25
		got := fieldConsole(t, Float32Ptr("key", &value))
		want := fieldConsole(t, Float32("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *float32
		if got := fieldConsole(t, Float32Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Float64", func(t *testing.T) {
		var value float64 = -2.5
		got := fieldConsole(t, Float64Ptr("key", &value))
		want := fieldConsole(t, Float64("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *float64
		if got := fieldConsole(t, Float64Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Complex128", func(t *testing.T) {
		var value complex128 = 1 + 2i
		got := fieldConsole(t, Complex128Ptr("key", &value))
		want := fieldConsole(t, Complex128("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *complex128
		if got := fieldConsole(t, Complex128Ptr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Binary", func(t *testing.T) {
		var value []byte = []byte{0, 255}
		got := fieldConsole(t, BinaryPtr("key", &value))
		want := fieldConsole(t, Binary("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *[]byte
		if got := fieldConsole(t, BinaryPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("ByteString", func(t *testing.T) {
		var value []byte = []byte("hello")
		got := fieldConsole(t, ByteStringPtr("key", &value))
		want := fieldConsole(t, ByteString("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *[]byte
		if got := fieldConsole(t, ByteStringPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("UnsafePointer", func(t *testing.T) {
		var value unsafe.Pointer = nil
		got := fieldConsole(t, UnsafePointerPtr("key", &value))
		want := fieldConsole(t, UnsafePointer("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *unsafe.Pointer
		if got := fieldConsole(t, UnsafePointerPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Array", func(t *testing.T) {
		var value []int = []int{1, 2}
		got := fieldConsole(t, ArrayPtr("key", &value))
		want := fieldConsole(t, Array("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *[]int
		if got := fieldConsole(t, ArrayPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Object", func(t *testing.T) {
		var value map[string]int = map[string]int{"a": 1}
		got := fieldConsole(t, ObjectPtr("key", &value))
		want := fieldConsole(t, Object("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *map[string]int
		if got := fieldConsole(t, ObjectPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Chan", func(t *testing.T) {
		var value chan int = make(chan int)
		got := fieldConsole(t, ChanPtr("key", &value))
		want := fieldConsole(t, Chan("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *chan int
		if got := fieldConsole(t, ChanPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Func", func(t *testing.T) {
		var value func() = func() {}
		got := fieldConsole(t, FuncPtr("key", &value))
		want := fieldConsole(t, Func("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing *func()
		if got := fieldConsole(t, FuncPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	t.Run("Pointer", func(t *testing.T) {
		var value *int = new(int)
		got := fieldConsole(t, PointerPtr("key", &value))
		want := fieldConsole(t, Pointer("key", value))
		if got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
		var missing **int
		if got := fieldConsole(t, PointerPtr("key", missing)); got != fieldConsole(t, Null("key")) {
			t.Fatalf("nil pointer: %q", got)
		}
	})
	value := 3
	field := IntPtr("key", &value)
	value = 4
	if got := fieldConsole(t, field); got != fieldConsole(t, Int("key", 3)) {
		t.Fatalf("pointer value not captured: %q", got)
	}
}
