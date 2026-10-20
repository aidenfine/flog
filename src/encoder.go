package src

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ObjectEncoder encodes named fields. Methods return serialization errors.
type ObjectEncoder interface {
	AddNull(string) error
	AddBool(string, bool) error
	AddString(string, string) error
	AddBinary(string, []byte) error
	AddByteString(string, []byte) error
	AddReflected(string, any) error
	AddArray(string, ArrayMarshaler) error
	AddObject(string, ObjectMarshaler) error
	AddInt(string, int) error
	AddInt8(string, int8) error
	AddInt16(string, int16) error
	AddInt32(string, int32) error
	AddInt64(string, int64) error
	AddUint(string, uint) error
	AddUint8(string, uint8) error
	AddUint16(string, uint16) error
	AddUint32(string, uint32) error
	AddUint64(string, uint64) error
	AddUintptr(string, uintptr) error
	AddFloat32(string, float32) error
	AddFloat64(string, float64) error
}

// ArrayEncoder encodes elements in order.
type ArrayEncoder interface {
	AppendReflected(any) error
	AppendArray(ArrayMarshaler) error
	AppendObject(ObjectMarshaler) error
	AppendBool(bool) error
	AppendString(string) error
	AppendInt(int) error
	AppendInt8(int8) error
	AppendInt16(int16) error
	AppendInt32(int32) error
	AppendInt64(int64) error
	AppendUint(uint) error
	AppendUint8(uint8) error
	AppendUint16(uint16) error
	AppendUint32(uint32) error
	AppendUint64(uint64) error
	AppendUintptr(uintptr) error
	AppendFloat32(float32) error
	AppendFloat64(float64) error
}

type ObjectMarshaler interface{ MarshalLogObject(ObjectEncoder) error }
type ArrayMarshaler interface{ MarshalLogArray(ArrayEncoder) error }

type jsonObjectEncoder struct {
	buf   bytes.Buffer
	count int
}

func (e *jsonObjectEncoder) add(key string, value []byte) error {

	k, err := json.Marshal(key)
	if err != nil {
		return err
	}
	if e.count > 0 {
		e.buf.WriteString(", ")
	}
	e.buf.Write(k)
	e.buf.WriteString(": ")
	e.buf.Write(value)
	e.count++
	return nil
}
func (e *jsonObjectEncoder) AddReflected(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return e.add(key, data)
}
func (e *jsonObjectEncoder) AddNull(key string) error             { return e.add(key, []byte("null")) }
func (e *jsonObjectEncoder) AddBool(key string, v bool) error     { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddString(key string, v string) error { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddBinary(key string, v []byte) error { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddByteString(key string, v []byte) error {
	return e.AddString(key, string(v))
}
func (e *jsonObjectEncoder) AddObject(key string, v ObjectMarshaler) error {
	child := &jsonObjectEncoder{}
	if err := v.MarshalLogObject(child); err != nil {
		return err
	}
	return e.add(key, []byte("{"+child.buf.String()+"}"))
}
func (e *jsonObjectEncoder) AddArray(key string, v ArrayMarshaler) error {
	child := &jsonArrayEncoder{}
	if err := v.MarshalLogArray(child); err != nil {
		return err
	}
	return e.add(key, []byte("["+child.buf.String()+"]"))
}
func (e *jsonObjectEncoder) AddInt(key string, v int) error         { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddInt8(key string, v int8) error       { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddInt16(key string, v int16) error     { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddInt32(key string, v int32) error     { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddInt64(key string, v int64) error     { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddUint(key string, v uint) error       { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddUint8(key string, v uint8) error     { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddUint16(key string, v uint16) error   { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddUint32(key string, v uint32) error   { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddUint64(key string, v uint64) error   { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddUintptr(key string, v uintptr) error { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddFloat32(key string, v float32) error { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddFloat64(key string, v float64) error { return e.AddReflected(key, v) }

type jsonArrayEncoder struct {
	buf   bytes.Buffer
	count int
}

func (e *jsonArrayEncoder) append(data []byte) error {
	if e.count > 0 {
		e.buf.WriteString(", ")
	}
	e.buf.Write(data)
	e.count++
	return nil
}
func (e *jsonArrayEncoder) AppendReflected(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return e.append(data)
}
func (e *jsonArrayEncoder) AppendObject(v ObjectMarshaler) error {
	child := &jsonObjectEncoder{}
	if err := v.MarshalLogObject(child); err != nil {
		return err
	}
	return e.append([]byte("{" + child.buf.String() + "}"))
}
func (e *jsonArrayEncoder) AppendArray(v ArrayMarshaler) error {
	child := &jsonArrayEncoder{}
	if err := v.MarshalLogArray(child); err != nil {
		return err
	}
	return e.append([]byte("[" + child.buf.String() + "]"))
}
func (e *jsonArrayEncoder) AppendBool(v bool) error       { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendString(v string) error   { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendInt(v int) error         { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendInt8(v int8) error       { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendInt16(v int16) error     { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendInt32(v int32) error     { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendInt64(v int64) error     { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendUint(v uint) error       { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendUint8(v uint8) error     { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendUint16(v uint16) error   { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendUint32(v uint32) error   { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendUint64(v uint64) error   { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendUintptr(v uintptr) error { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendFloat32(v float32) error { return e.AppendReflected(v) }
func (e *jsonArrayEncoder) AppendFloat64(v float64) error { return e.AppendReflected(v) }

func encodeFields(fields []Field) ([]byte, error) {
	enc := &jsonObjectEncoder{}
	for i, field := range fields {
		// Positional arguments retain their order and receive a stable name.
		if field.Key == "" {
			field.Key = fmt.Sprintf("arg%d", i)
		}
		if err := field.AddTo(enc); err != nil {
			return nil, fmt.Errorf("field %q: %w", field.Key, err)
		}
	}
	return []byte("{" + enc.buf.String() + "}"), nil
}
