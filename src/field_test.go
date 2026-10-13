package src

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestMakeFields(t *testing.T) {
	type namedInt int32
	type namedString string
	type namedBytes []byte
	number := 42
	fn := func() {}
	ch := make(chan int)

	tests := []struct {
		name    string
		value   any
		typ     FieldType
		integer int64
		str     string
	}{
		{name: "nil", typ: NoneType},
		{name: "true", value: true, typ: BoolType, integer: 1},
		{name: "false", value: false, typ: BoolType},
		{name: "int", value: int(-42), typ: IntType, integer: -42},
		{name: "int8", value: int8(-128), typ: Int8Type, integer: -128},
		{name: "int16", value: int16(-32768), typ: Int16Type, integer: -32768},
		{name: "int32", value: int32(-2147483648), typ: Int32Type, integer: -2147483648},
		{name: "int64", value: int64(-9223372036854775808), typ: Int64Type, integer: -9223372036854775808},
		{name: "uint", value: uint(42), typ: UintType},
		{name: "uint8", value: uint8(255), typ: Uint8Type},
		{name: "uint16", value: uint16(65535), typ: Uint16Type},
		{name: "uint32", value: uint32(4294967295), typ: Uint32Type},
		{name: "uint64", value: uint64(18446744073709551615), typ: Uint64Type},
		{name: "uintptr", value: uintptr(42), typ: UintptrType},
		{name: "byte alias", value: byte(255), typ: Uint8Type},
		{name: "rune alias", value: rune('世'), typ: Int32Type, integer: '世'},
		{name: "float32", value: float32(1.25), typ: Float32Type},
		{name: "float64", value: float64(-2.5), typ: Float64Type},
		{name: "complex64", value: complex64(1 + 2i), typ: Complex128Type},
		{name: "complex128", value: complex128(3 - 4i), typ: Complex128Type},
		{name: "string", value: "hello", typ: StringType, str: "hello"},
		{name: "empty string", value: "", typ: StringType},
		{name: "array", value: [2]int{1, 2}, typ: ArrayType},
		{name: "byte array", value: [2]byte{1, 2}, typ: ArrayType},
		{name: "slice", value: []int{1, 2}, typ: ArrayType},
		{name: "empty slice", value: []int{}, typ: ArrayType},
		{name: "nil slice", value: []int(nil), typ: ArrayType},
		{name: "binary", value: []byte{0, 255}, typ: BinaryType},
		{name: "nil binary", value: []byte(nil), typ: BinaryType},
		{name: "map", value: map[string]int{"answer": 42}, typ: ObjectType},
		{name: "nil map", value: map[string]int(nil), typ: ObjectType},
		{name: "struct", value: struct{ Count int }{42}, typ: ObjectType},
		{name: "pointer", value: &number, typ: PointerType},
		{name: "nil pointer", value: (*int)(nil), typ: PointerType},
		{name: "channel", value: ch, typ: ChanType},
		{name: "nil channel", value: (chan int)(nil), typ: ChanType},
		{name: "function", value: fn, typ: FuncType},
		{name: "nil function", value: (func())(nil), typ: FuncType},
		{name: "unsafe pointer", value: unsafe.Pointer(&number), typ: UnsafePointerType},
		{name: "nil unsafe pointer", value: unsafe.Pointer(nil), typ: UnsafePointerType},
		{name: "named int", value: namedInt(-42), typ: Int32Type, integer: -42},
		{name: "named string", value: namedString("hello"), typ: StringType, str: "hello"},
		{name: "named bytes", value: namedBytes{1, 2}, typ: BinaryType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := MakeFields(tt.value)
			if len(fields) != 1 {
				t.Fatalf("got %d fields, want 1", len(fields))
			}
			got := fields[0]
			if got.Key != "" || got.Type != tt.typ || got.Integer != tt.integer || got.String != tt.str {
				t.Errorf("got key=%q type=%d integer=%d string=%q; want key=\"\" type=%d integer=%d string=%q",
					got.Key, got.Type, got.Integer, got.String, tt.typ, tt.integer, tt.str)
			}
			// type and code pointer separately.
			if tt.typ == FuncType {
				if reflect.TypeOf(got.Any) != reflect.TypeOf(tt.value) ||
					reflect.ValueOf(got.Any).Pointer() != reflect.ValueOf(tt.value).Pointer() {
					t.Error("original function was not preserved")
				}
			} else if !reflect.DeepEqual(got.Any, tt.value) {
				t.Errorf("Any = %#v (%T), want %#v (%T)", got.Any, got.Any, tt.value, tt.value)
			}
		})
	}
}

func TestMakeFieldsArguments(t *testing.T) {
	tests := []struct {
		name   string
		values []any
		want   []Field
	}{
		{name: "no arguments", want: []Field{}},
		{
			name:   "multiple arguments preserve order",
			values: []any{int32(7), "hello", nil, true},
			want: []Field{
				{Type: Int32Type, Integer: 7, Any: int32(7)},
				{Type: StringType, String: "hello", Any: "hello"},
				{Type: NoneType},
				{Type: BoolType, Integer: 1, Any: true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeFields(tt.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeFields() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
