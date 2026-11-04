package flog

import (
	"reflect"
	"slices"
	"unsafe"

	"github.com/aidenfine/flog/src"
)

// Field is a named value added to a log entry.
type Field = src.Field
type ObjectEncoder = src.ObjectEncoder
type ArrayEncoder = src.ArrayEncoder
type ObjectMarshaler = src.ObjectMarshaler
type ArrayMarshaler = src.ArrayMarshaler

// Any creates a named field, supporting custom marshalers and arbitrary values.
func Any(key string, value any) Field {
	field := src.MakeFields(value)[0]
	field.Key = key
	return field
}

// Frontend for strict logging.

// Constructs a String field.
func String(key, value string) Field {
	return Field{
		Key:    key,
		Type:   src.StringType,
		String: value,
	}
}

// Constructs a Bool field.
func Bool(key string, value bool) Field {
	var v int64
	// 1- true
	if value {
		v = 1
	}
	return Field{
		Key:     key,
		Type:    src.BoolType,
		Integer: v,
	}
}

// Constructs a Binary field.
func Binary(key string, value []byte) Field {
	return Field{Key: key, Type: src.BinaryType, Any: value}
}

// Constructs a ByteString field.
func ByteString(key string, value []byte) Field {
	return Field{Key: key, Type: src.ByteStringType, Any: value}
}

// Construct Int field.
func Int(key string, value int) Field {
	return Field{Key: key, Type: src.IntType, Integer: int64(value)}
}

func Int8(key string, value int8) Field {
	return Field{Key: key, Type: src.Int8Type, Integer: int64(value)}
}

func Int16(key string, value int16) Field {
	return Field{Key: key, Type: src.Int16Type, Integer: int64(value)}
}

func Int32(key string, value int32) Field {
	return Field{Key: key, Type: src.Int32Type, Integer: int64(value)}
}

// Construct Int64 field.
func Int64(key string, value int64) Field {
	return Field{Key: key, Type: src.Int64Type, Integer: value}
}

// Uint constructors
func Uint(key string, value uint) Field {
	return Field{Key: key, Type: src.UintType, Any: value}
}
func Uint8(key string, value uint8) Field {
	return Field{Key: key, Type: src.Uint8Type, Any: value}
}

func Uint16(key string, value uint16) Field {
	return Field{Key: key, Type: src.Uint16Type, Any: value}
}

func Uint32(key string, value uint32) Field {
	return Field{Key: key, Type: src.Uint32Type, Any: value}
}

func Uint64(key string, value uint64) Field {
	return Field{Key: key, Type: src.Uint64Type, Any: value}
}

// Unit16 is retained for compatibility. Use Uint16 instead.
func Unit16(key string, value uint16) Field { return Uint16(key, value) }

// Null creates an explicitly null field.
func Null(key string) Field { return Field{Key: key, Type: src.NoneType} }

// Invalid creates a field that reports an encoding error when written.
func Invalid(key string, value any) Field { return Field{Key: key, Type: src.Invalidtype, Any: value} }

// Array accepts a Go array, slice, or custom ArrayMarshaler.
func Array(key string, value any) Field {
	if _, ok := value.(ArrayMarshaler); ok {
		return Field{Key: key, Type: src.ArrayType, Any: value}
	}
	return kindField(key, value, src.ArrayType, reflect.Array, reflect.Slice)
}

// Object accepts a Go map, struct, or custom ObjectMarshaler.
func Object(key string, value any) Field {
	if _, ok := value.(ObjectMarshaler); ok {
		return Field{Key: key, Type: src.ObjectType, Any: value}
	}
	return kindField(key, value, src.ObjectType, reflect.Map, reflect.Struct)
}

// Pointer retains a raw Go pointer for reflected encoding.
func Pointer[T any](key string, value *T) Field {
	return Field{Key: key, Type: src.PointerType, Any: value}
}

// Chan accepts bidirectional, send-only, and receive-only channels.
func Chan(key string, value any) Field { return kindField(key, value, src.ChanType, reflect.Chan) }

// Func accepts functions of any signature.
func Func(key string, value any) Field { return kindField(key, value, src.FuncType, reflect.Func) }

// UnsafePointer creates a field containing a raw unsafe pointer.
func UnsafePointer(key string, value unsafe.Pointer) Field {
	return Field{Key: key, Type: src.UnsafePointerType, Any: value}
}

func kindField(key string, value any, typ src.FieldType, kinds ...reflect.Kind) Field {
	if value != nil && slices.Contains(kinds, reflect.TypeOf(value).Kind()) {
		return Field{Key: key, Type: typ, Any: value}
	}
	return Invalid(key, value)
}

// Uintptr creates a named uintptr field.
func Uintptr(key string, value uintptr) Field {
	return Field{Key: key, Type: src.UintptrType, Any: value}
}

// Float32 creates a named float32 field.
func Float32(key string, value float32) Field {
	return Field{Key: key, Type: src.Float32Type, Any: value}
}

// Float64 creates a named float64 field.
func Float64(key string, value float64) Field {
	return Field{Key: key, Type: src.Float64Type, Any: value}
}

// Complex128 creates a named complex128 field.
func Complex128(key string, value complex128) Field {
	return Field{Key: key, Type: src.Complex128Type, Any: value}
}

// StringPtr logs a snapshot of the ptr value, or null for nil.
func StringPtr(key string, value *string) Field {
	if value == nil {
		return Null(key)
	}
	return String(key, *value)
}

// BoolPtr logs a snapshot of the ptr value, or null for nil.
func BoolPtr(key string, value *bool) Field {
	if value == nil {
		return Null(key)
	}
	return Bool(key, *value)
}

// IntPtr logs a snapshot of the ptr value, or null for nil.
func IntPtr(key string, value *int) Field {
	if value == nil {
		return Null(key)
	}
	return Int(key, *value)
}

// Int8Ptr logs a snapshot of the ptr value, or null for nil.
func Int8Ptr(key string, value *int8) Field {
	if value == nil {
		return Null(key)
	}
	return Int8(key, *value)
}

// Int16Ptr logs a snapshot of the ptr value, or null for nil.
func Int16Ptr(key string, value *int16) Field {
	if value == nil {
		return Null(key)
	}
	return Int16(key, *value)
}

// Int32Ptr logs a snapshot of the ptr value, or null for nil.
func Int32Ptr(key string, value *int32) Field {
	if value == nil {
		return Null(key)
	}
	return Int32(key, *value)
}

// Int64Ptr logs a snapshot of the ptr value, or null for nil.
func Int64Ptr(key string, value *int64) Field {
	if value == nil {
		return Null(key)
	}
	return Int64(key, *value)
}

// UintPtr logs a snapshot of the ptr value, or null for nil.
func UintPtr(key string, value *uint) Field {
	if value == nil {
		return Null(key)
	}
	return Uint(key, *value)
}

// Uint8Ptr logs a snapshot of the ptr value, or null for nil.
func Uint8Ptr(key string, value *uint8) Field {
	if value == nil {
		return Null(key)
	}
	return Uint8(key, *value)
}

// Uint16Ptr logs a snapshot of the ptr value, or null for nil.
func Uint16Ptr(key string, value *uint16) Field {
	if value == nil {
		return Null(key)
	}
	return Uint16(key, *value)
}

// Uint32Ptr logs a snapshot of the ptr value, or null for nil.
func Uint32Ptr(key string, value *uint32) Field {
	if value == nil {
		return Null(key)
	}
	return Uint32(key, *value)
}

// Uint64Ptr logs a snapshot of the ptr value, or null for nil.
func Uint64Ptr(key string, value *uint64) Field {
	if value == nil {
		return Null(key)
	}
	return Uint64(key, *value)
}

// UintptrPtr logs a snapshot of the ptr value, or null for nil.
func UintptrPtr(key string, value *uintptr) Field {
	if value == nil {
		return Null(key)
	}
	return Uintptr(key, *value)
}

// Float32Ptr logs a snapshot of the ptr value, or null for nil.
func Float32Ptr(key string, value *float32) Field {
	if value == nil {
		return Null(key)
	}
	return Float32(key, *value)
}

// Float64Ptr logs a snapshot of the ptr value, or null for nil.
func Float64Ptr(key string, value *float64) Field {
	if value == nil {
		return Null(key)
	}
	return Float64(key, *value)
}

// Complex128Ptr logs a snapshot of the ptr value, or null for nil.
func Complex128Ptr(key string, value *complex128) Field {
	if value == nil {
		return Null(key)
	}
	return Complex128(key, *value)
}

// BinaryPtr logs a snapshot of the ptr value, or null for nil.
func BinaryPtr(key string, value *[]byte) Field {
	if value == nil {
		return Null(key)
	}
	return Binary(key, *value)
}

// ByteStringPtr logs a snapshot of the ptr value, or null for nil.
func ByteStringPtr(key string, value *[]byte) Field {
	if value == nil {
		return Null(key)
	}
	return ByteString(key, *value)
}

// UnsafePointerPtr logs a snapshot of the ptr value, or null for nil.
func UnsafePointerPtr(key string, value *unsafe.Pointer) Field {
	if value == nil {
		return Null(key)
	}
	return UnsafePointer(key, *value)
}

// ArrayPtr logs the ptr array or slice, or null for nil.
func ArrayPtr[T any](key string, value *T) Field {
	if value == nil {
		return Null(key)
	}
	return Array(key, *value)
}

// ObjectPtr logs the ptr object, or null for nil.
func ObjectPtr[T any](key string, value *T) Field {
	if value == nil {
		return Null(key)
	}
	return Object(key, *value)
}

// ChanPtr logs the ptr channel, or null for nil.
func ChanPtr[T any](key string, value *T) Field {
	if value == nil {
		return Null(key)
	}
	return Chan(key, *value)
}

// FuncPtr logs the ptr function, or null for nil.
func FuncPtr[T any](key string, value *T) Field {
	if value == nil {
		return Null(key)
	}
	return Func(key, *value)
}

// PointerPtr logs the raw pointer held by value, or null for nil.
func PointerPtr[T any](key string, value **T) Field {
	if value == nil {
		return Null(key)
	}
	return Pointer(key, *value)
}
