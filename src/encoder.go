package src

import (
	"time"
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

func EpochTimeEncoder(t time.Time) float64 {
	nanoSeconds := t.UnixNano()
	return float64(nanoSeconds) / float64(time.Second)
}

func EpochMillisTimeEncoder(t time.Time) float64 {
	nanos := t.UnixNano()
	return float64(nanos) / float64(time.Millisecond)
}

func EpochNanosTimeEncoder(t time.Time) int64 {
	return t.UnixNano()
}

func encodeTime(t time.Time, format string) string {
	return t.Format(format)
}

func RFC3339TimeEncoder(t time.Time) string {
	return encodeTime(t, time.RFC3339)
}

func RFC3339NanoEncoder(t time.Time) string {
	return encodeTime(t, time.RFC3339Nano)
}

// bc why not!
func KitchenEncoder(t time.Time) string {
	return encodeTime(t, time.Kitchen)
}
