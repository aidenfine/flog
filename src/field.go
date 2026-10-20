package src

import (
	"fmt"
	"reflect"
)

type FieldType uint8

const (
	NoneType FieldType = iota
	ArrayType
	ObjectType
	BinaryType
	BoolType
	ByteStringType
	Int64Type
	Uint64Type
	Float64Type
	Complex128Type
	StringType
	PointerType
	ChanType
	FuncType
	UnsafePointerType
	IntType
	Int8Type
	Int16Type
	Int32Type
	UintType
	Uint8Type
	Uint16Type
	Uint32Type
	UintptrType
	Float32Type
)

type Field struct {
	Key     string
	Type    FieldType
	Integer int64
	String  string
	Any     any
}

func AppendFields(dst []Field, values ...any) []Field {
	for _, value := range values {
		if field, ok := value.(Field); ok {
			dst = append(dst, field)
			continue
		}
		field := Field{Any: value}
		if _, ok := value.(ObjectMarshaler); ok {
			field.Type = ObjectType
			dst = append(dst, field)
			continue
		}
		if _, ok := value.(ArrayMarshaler); ok {
			field.Type = ArrayType
			dst = append(dst, field)
			continue
		}
		v := reflect.ValueOf(value)
		if !v.IsValid() {
			dst = append(dst, field)
			continue
		}
		switch v.Kind() {
		case reflect.Bool:
			field.Type = BoolType
			if v.Bool() {
				field.Integer = 1
			}
		case reflect.Int:
			field.Type = IntType
			field.Integer = v.Int()
		case reflect.Int8:
			field.Type = Int8Type
			field.Integer = v.Int()
		case reflect.Int16:
			field.Type = Int16Type
			field.Integer = v.Int()
		case reflect.Int32:
			field.Type = Int32Type
			field.Integer = v.Int()
		case reflect.Int64:
			field.Type = Int64Type
			field.Integer = v.Int()
		case reflect.Uint:
			field.Type = UintType
		case reflect.Uint8:
			field.Type = Uint8Type
		case reflect.Uint16:
			field.Type = Uint16Type
		case reflect.Uint32:
			field.Type = Uint32Type
		case reflect.Uint64:
			// Preserve unsigned values in Any to avoid overflowing int64.
			field.Type = Uint64Type
		case reflect.Uintptr:
			field.Type = UintptrType
		case reflect.Float32:
			field.Type = Float32Type
		case reflect.Float64:
			field.Type = Float64Type
		case reflect.Complex64, reflect.Complex128:
			field.Type = Complex128Type
		case reflect.String:
			field.Type = StringType
			field.String = v.String()
		case reflect.Array:
			field.Type = ArrayType
		case reflect.Slice:
			if v.Type().Elem().Kind() == reflect.Uint8 {
				field.Type = BinaryType
			} else {
				field.Type = ArrayType
			}
		case reflect.Map, reflect.Struct:
			field.Type = ObjectType
		case reflect.Pointer:
			field.Type = PointerType
		case reflect.Chan:
			field.Type = ChanType
		case reflect.Func:
			field.Type = FuncType
		case reflect.UnsafePointer:
			field.Type = UnsafePointerType
		}
		dst = append(dst, field)

	}
	return dst
}

func MakeFields(fields ...any) []Field {
	return AppendFields(nil, fields...)
}

// AddTo dispatches a field to the encoder without choosing an output format.
func (f Field) AddTo(enc ObjectEncoder) error {
	switch f.Type {
	case NoneType:
		return enc.AddNull(f.Key)
	case BoolType:
		return enc.AddBool(f.Key, f.Integer != 0)
	case StringType:
		return enc.AddString(f.Key, f.String)
	case IntType:
		return enc.AddInt(f.Key, int(f.Integer))
	case Int8Type:
		return enc.AddInt8(f.Key, int8(f.Integer))
	case Int16Type:
		return enc.AddInt16(f.Key, int16(f.Integer))
	case Int32Type:
		return enc.AddInt32(f.Key, int32(f.Integer))
	case Int64Type:
		return enc.AddInt64(f.Key, int64(f.Integer))
	case BinaryType, ByteStringType:
		if f.Type == ByteStringType && f.Any == nil {
			return enc.AddString(f.Key, f.String)
		}
		v := reflect.ValueOf(f.Any)
		if !v.IsValid() || v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("expected byte slice, got %T", f.Any)
		}
		if f.Type == BinaryType {
			return enc.AddBinary(f.Key, v.Bytes())
		}
		return enc.AddByteString(f.Key, v.Bytes())
	case ArrayType:
		if v, ok := f.Any.(ArrayMarshaler); ok {
			return enc.AddArray(f.Key, v)
		}
		return enc.AddReflected(f.Key, f.Any)
	case ObjectType:
		if v, ok := f.Any.(ObjectMarshaler); ok {
			return enc.AddObject(f.Key, v)
		}
		return enc.AddReflected(f.Key, f.Any)
	case UintType:
		v := reflect.ValueOf(f.Any)
		if !v.IsValid() || v.Kind() != reflect.Uint {
			return fmt.Errorf("expected uint, got %T", f.Any)
		}
		return enc.AddUint(f.Key, uint(v.Uint()))
	case Uint8Type:
		v := reflect.ValueOf(f.Any)
		if !v.IsValid() || v.Kind() != reflect.Uint8 {
			return fmt.Errorf("expected uint8, got %T", f.Any)
		}
		return enc.AddUint8(f.Key, uint8(v.Uint()))
	case Uint16Type:
		v := reflect.ValueOf(f.Any)
		if !v.IsValid() || v.Kind() != reflect.Uint16 {
			return fmt.Errorf("expected uint16, got %T", f.Any)
		}
		return enc.AddUint16(f.Key, uint16(v.Uint()))
	case Uint32Type:
		v := reflect.ValueOf(f.Any)
		if !v.IsValid() || v.Kind() != reflect.Uint32 {
			return fmt.Errorf("expected uint32, got %T", f.Any)
		}
		return enc.AddUint32(f.Key, uint32(v.Uint()))
	case Uint64Type:
		v := reflect.ValueOf(f.Any)
		if !v.IsValid() || v.Kind() != reflect.Uint64 {
			return fmt.Errorf("expected uint64, got %T", f.Any)
		}
		return enc.AddUint64(f.Key, uint64(v.Uint()))
	case UintptrType:
		v := reflect.ValueOf(f.Any)
		if !v.IsValid() || v.Kind() != reflect.Uintptr {
			return fmt.Errorf("expected uintptr, got %T", f.Any)
		}
		return enc.AddUintptr(f.Key, uintptr(v.Uint()))
	case Float32Type:
		v := reflect.ValueOf(f.Any)
		if !v.IsValid() || v.Kind() != reflect.Float32 {
			return fmt.Errorf("expected float32, got %T", f.Any)
		}
		return enc.AddFloat32(f.Key, float32(v.Float()))
	case Float64Type:
		v := reflect.ValueOf(f.Any)
		if !v.IsValid() || v.Kind() != reflect.Float64 {
			return fmt.Errorf("expected float64, got %T", f.Any)
		}
		return enc.AddFloat64(f.Key, float64(v.Float()))
	case Complex128Type:
		v := reflect.ValueOf(f.Any)
		if !v.IsValid() || (v.Kind() != reflect.Complex64 && v.Kind() != reflect.Complex128) {
			return fmt.Errorf("expected complex value, got %T", f.Any)
		}
		return enc.AddString(f.Key, fmt.Sprint(f.Any))
	case PointerType, ChanType, FuncType, UnsafePointerType:
		return enc.AddReflected(f.Key, f.Any)
	default:
		return fmt.Errorf("unknown field type %d", f.Type)
	}
}
