package src

import (
	"bytes"
	"fmt"
	"strconv"
)

// consoleEncoder writes values directly into the log buffer. Only nested
// objects include keys; top-level fields retain positional console output.
type consoleEncoder struct {
	buf   *bytes.Buffer
	count int
	named bool
}

func (e *consoleEncoder) next(key string) {
	if e.count > 0 {
		e.buf.WriteByte(' ')
	}
	e.count++
	if e.named {
		e.buf.WriteString(key)
		e.buf.WriteByte(':')
	}
}

func (e *consoleEncoder) AddNull(key string) error {
	e.next(key)
	e.buf.WriteString("<nil>")
	return nil
}
func (e *consoleEncoder) AddReflected(key string, value any) error {
	e.next(key)
	fmt.Fprint(e.buf, value)
	return nil
}
func (e *consoleEncoder) AddString(key, value string) error {
	e.next(key)
	e.buf.WriteString(value)
	return nil
}
func (e *consoleEncoder) AddBinary(key string, value []byte) error {
	return e.AddReflected(key, value)
}
func (e *consoleEncoder) AddByteString(key string, value []byte) error {
	e.next(key)
	e.buf.Write(value)
	return nil
}
func (e *consoleEncoder) AddObject(key string, value ObjectMarshaler) error {
	e.next(key)
	e.buf.WriteByte('{')
	child := &consoleEncoder{buf: e.buf, named: true}
	if err := value.MarshalLogObject(child); err != nil {
		return err
	}
	e.buf.WriteByte('}')
	return nil
}
func (e *consoleEncoder) AddArray(key string, value ArrayMarshaler) error {
	e.next(key)
	e.buf.WriteByte('[')
	child := &consoleEncoder{buf: e.buf}
	if err := value.MarshalLogArray(child); err != nil {
		return err
	}
	e.buf.WriteByte(']')
	return nil
}

func (e *consoleEncoder) AddBool(key string, value bool) error {
	e.next(key)
	e.buf.WriteString(strconv.FormatBool(value))
	return nil
}

func (e *consoleEncoder) AddInt(key string, value int) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendInt(scratch[:0], int64(value), 10))
	return nil
}

func (e *consoleEncoder) AddInt8(key string, value int8) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendInt(scratch[:0], int64(value), 10))
	return nil
}

func (e *consoleEncoder) AddInt16(key string, value int16) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendInt(scratch[:0], int64(value), 10))
	return nil
}

func (e *consoleEncoder) AddInt32(key string, value int32) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendInt(scratch[:0], int64(value), 10))
	return nil
}

func (e *consoleEncoder) AddInt64(key string, value int64) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendInt(scratch[:0], int64(value), 10))
	return nil
}

func (e *consoleEncoder) AddUint(key string, value uint) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendUint(scratch[:0], uint64(value), 10))
	return nil
}

func (e *consoleEncoder) AddUint8(key string, value uint8) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendUint(scratch[:0], uint64(value), 10))
	return nil
}

func (e *consoleEncoder) AddUint16(key string, value uint16) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendUint(scratch[:0], uint64(value), 10))
	return nil
}

func (e *consoleEncoder) AddUint32(key string, value uint32) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendUint(scratch[:0], uint64(value), 10))
	return nil
}

func (e *consoleEncoder) AddUint64(key string, value uint64) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendUint(scratch[:0], uint64(value), 10))
	return nil
}

func (e *consoleEncoder) AddUintptr(key string, value uintptr) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendUint(scratch[:0], uint64(value), 10))
	return nil
}

func (e *consoleEncoder) AddFloat32(key string, value float32) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendFloat(scratch[:0], float64(value), 'g', -1, 32))
	return nil
}

func (e *consoleEncoder) AddFloat64(key string, value float64) error {
	e.next(key)
	var scratch [64]byte
	e.buf.Write(strconv.AppendFloat(scratch[:0], float64(value), 'g', -1, 64))
	return nil
}

func (e *consoleEncoder) AppendReflected(value any) error { return e.AddReflected("", value) }

func (e *consoleEncoder) AppendArray(value ArrayMarshaler) error { return e.AddArray("", value) }

func (e *consoleEncoder) AppendObject(value ObjectMarshaler) error { return e.AddObject("", value) }

func (e *consoleEncoder) AppendBool(value bool) error { return e.AddBool("", value) }

func (e *consoleEncoder) AppendString(value string) error { return e.AddString("", value) }

func (e *consoleEncoder) AppendInt(value int) error { return e.AddInt("", value) }

func (e *consoleEncoder) AppendInt8(value int8) error { return e.AddInt8("", value) }

func (e *consoleEncoder) AppendInt16(value int16) error { return e.AddInt16("", value) }

func (e *consoleEncoder) AppendInt32(value int32) error { return e.AddInt32("", value) }

func (e *consoleEncoder) AppendInt64(value int64) error { return e.AddInt64("", value) }

func (e *consoleEncoder) AppendUint(value uint) error { return e.AddUint("", value) }

func (e *consoleEncoder) AppendUint8(value uint8) error { return e.AddUint8("", value) }

func (e *consoleEncoder) AppendUint16(value uint16) error { return e.AddUint16("", value) }

func (e *consoleEncoder) AppendUint32(value uint32) error { return e.AddUint32("", value) }

func (e *consoleEncoder) AppendUint64(value uint64) error { return e.AddUint64("", value) }

func (e *consoleEncoder) AppendUintptr(value uintptr) error { return e.AddUintptr("", value) }

func (e *consoleEncoder) AppendFloat32(value float32) error { return e.AddFloat32("", value) }

func (e *consoleEncoder) AppendFloat64(value float64) error { return e.AddFloat64("", value) }

func writeConsoleValues(buf *bytes.Buffer, fields []Field) error {
	enc := &consoleEncoder{buf: buf}
	for i, field := range fields {
		if err := field.AddTo(enc); err != nil {
			key := field.Key
			if key == "" {
				key = fmt.Sprintf("arg%d", i)
			}
			return fmt.Errorf("field %q: %w", key, err)
		}
	}
	return nil
}
