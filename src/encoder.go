package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
	"unicode/utf8"
)

// consoleEncoder writes console values, with JSON-like formatting for named fields.
type consoleEncoder struct {
	buf        *bytes.Buffer
	count      int
	named      bool
	structured bool
}

// Writer reuses a buffer across log entries.
type Writer struct {
	buf        bytes.Buffer
	consoleEnc consoleEncoder
	jsonEnc    jsonObjectEncoder
}

type ObjectMarshaler interface{ MarshalLogObject(ObjectEncoder) error }
type ArrayMarshaler interface{ MarshalLogArray(ArrayEncoder) error }

type jsonObjectEncoder struct {
	buf   *bytes.Buffer
	count int
}

// quote appends into the buffer's spare capacity without a temporary string.
func (e *consoleEncoder) quote(value string) {
	e.buf.Grow(4*len(value) + 2)
	e.buf.Write(strconv.AppendQuote(e.buf.AvailableBuffer(), value))
}

func (e *consoleEncoder) next(key string) {
	if e.count > 0 {
		if e.structured {
			e.buf.WriteString(", ")
		} else {
			e.buf.WriteByte(' ')
		}
	}
	e.count++
	if e.named {
		if e.structured {
			e.quote(key)
			e.buf.WriteString(": ")
		} else {
			e.buf.WriteString(key)
			e.buf.WriteByte(':')
		}
	} else if key != "" {
		e.buf.WriteString(key)
		e.buf.WriteByte('=')
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
	if e.structured {
		e.quote(value)
	} else {
		e.buf.WriteString(value)
	}
	return nil
}
func (e *consoleEncoder) AddBinary(key string, value []byte) error {
	return e.AddReflected(key, value)
}
func (e *consoleEncoder) AddByteString(key string, value []byte) error {
	return e.AddString(key, string(value))
}
func (e *consoleEncoder) AddObject(key string, value ObjectMarshaler) error {
	e.next(key)
	e.buf.WriteByte('{')
	child := &consoleEncoder{buf: e.buf, named: true, structured: e.structured}
	if err := value.MarshalLogObject(child); err != nil {
		return err
	}
	e.buf.WriteByte('}')
	return nil
}
func (e *consoleEncoder) AddArray(key string, value ArrayMarshaler) error {
	e.next(key)
	e.buf.WriteByte('[')
	child := &consoleEncoder{buf: e.buf, structured: e.structured}
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

func (w *Writer) writeConsoleValues(fields []Field) error {
	structured := false
	for _, field := range fields {
		if field.Key != "" {
			structured = true
			break
		}
	}

	w.consoleEnc = consoleEncoder{
		buf:        &w.buf,
		named:      structured,
		structured: structured,
	}

	if structured {
		w.buf.WriteByte('{')
	}
	for i, field := range fields {
		if duration, ok := field.Any.(time.Duration); ok {
			field.Type = StringType
			field.String = duration.String()
		}
		if err := field.AddTo(&w.consoleEnc); err != nil {
			key := field.Key
			if key == "" {
				key = fmt.Sprintf("arg%d", i)
			}
			return fmt.Errorf("field %q: %w", key, err)
		}
	}
	if structured {
		w.buf.WriteByte('}')
	}
	return nil
}

// ObjectEncoder and ArrayEncoder are shared by both output formats.

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

// writeJSONString matches encoding/json's string escaping. Avoiding using marshal seems to boost perf.
func writeJSONString(buf *bytes.Buffer, value string) {
	const hex = "0123456789abcdef"
	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(value); {
		c := value[i]
		if c < utf8.RuneSelf {
			if c >= 0x20 && c != '"' && c != '\\' && c != '<' && c != '>' && c != '&' {
				i++
				continue
			}
			buf.WriteString(value[start:i])
			switch c {
			case '"', '\\':
				buf.WriteByte('\\')
				buf.WriteByte(c)
			case '\n':
				buf.WriteString(`\n`)
			case '\r':
				buf.WriteString(`\r`)
			case '\t':
				buf.WriteString(`\t`)
			case '\b':
				buf.WriteString(`\b`)
			case '\f':
				buf.WriteString(`\f`)
			default:
				buf.WriteString(`\u00`)
				buf.WriteByte(hex[c>>4])
				buf.WriteByte(hex[c&15])
			}
			i++
			start = i
			continue
		}
		r, size := utf8.DecodeRuneInString(value[i:])
		if r == utf8.RuneError && size == 1 {
			buf.WriteString(value[start:i])
			buf.WriteString(`\ufffd`)
			i++
			start = i
		} else if r == '\u2028' || r == '\u2029' {
			buf.WriteString(value[start:i])
			buf.WriteString(`\u202`)
			buf.WriteByte(hex[r&15])
			i += size
			start = i
		} else {
			i += size
		}
	}
	buf.WriteString(value[start:])
	buf.WriteByte('"')
}

func (e *jsonObjectEncoder) key(key string) {
	if e.count > 0 {
		e.buf.WriteString(", ")
	}
	writeJSONString(e.buf, key)
	e.buf.WriteString(": ")
	e.count++
}
func (e *jsonObjectEncoder) add(key string, value []byte) error {
	e.key(key)
	e.buf.Write(value)
	return nil
}
func (e *jsonObjectEncoder) AddReflected(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return e.add(key, data)
}
func (e *jsonObjectEncoder) AddNull(key string) error { return e.add(key, []byte("null")) }
func (e *jsonObjectEncoder) AddBool(key string, v bool) error {
	e.key(key)
	if v {
		e.buf.WriteString("true")
	} else {
		e.buf.WriteString("false")
	}
	return nil
}
func (e *jsonObjectEncoder) AddString(key string, v string) error {
	e.key(key)
	writeJSONString(e.buf, v)
	return nil
}
func (e *jsonObjectEncoder) AddBinary(key string, v []byte) error { return e.AddReflected(key, v) }
func (e *jsonObjectEncoder) AddByteString(key string, v []byte) error {
	return e.AddString(key, string(v))
}
func (e *jsonObjectEncoder) AddObject(key string, v ObjectMarshaler) error {
	start, count := e.buf.Len(), e.count
	e.key(key)
	e.buf.WriteByte('{')
	child := jsonObjectEncoder{buf: e.buf}
	if err := v.MarshalLogObject(&child); err != nil {
		e.buf.Truncate(start)
		e.count = count
		return err
	}
	e.buf.WriteByte('}')
	return nil
}
func (e *jsonObjectEncoder) AddArray(key string, v ArrayMarshaler) error {
	start, count := e.buf.Len(), e.count
	e.key(key)
	e.buf.WriteByte('[')
	child := jsonArrayEncoder{buf: e.buf}
	if err := v.MarshalLogArray(&child); err != nil {
		e.buf.Truncate(start)
		e.count = count
		return err
	}
	e.buf.WriteByte(']')
	return nil
}
func (e *jsonObjectEncoder) AddInt64(key string, v int64) error {
	e.key(key)
	e.buf.Write(strconv.AppendInt(e.buf.AvailableBuffer(), v, 10))
	return nil
}
func (e *jsonObjectEncoder) AddUint64(key string, v uint64) error {
	e.key(key)
	e.buf.Write(strconv.AppendUint(e.buf.AvailableBuffer(), v, 10))
	return nil
}

// JSON uses decimal notation for [1e-6, 1e21) and removes leading exponent zeros.
func jsonFloat(v float64, bits int, dst []byte) ([]byte, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		value := any(v)
		if bits == 32 {
			value = float32(v)
		}
		return nil, &json.UnsupportedValueError{Value: reflect.ValueOf(value), Str: strconv.FormatFloat(v, 'g', -1, bits)}
	}
	format := byte('f')
	abs := math.Abs(v)
	if bits == 32 {
		abs = float64(float32(abs))
	}
	if abs != 0 && ((bits == 64 && (abs < 1e-6 || abs >= 1e21)) || (bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21))) {
		format = 'e'
	}
	dst = strconv.AppendFloat(dst, v, format, -1, bits)
	n := len(dst)
	if format == 'e' && n >= 4 && dst[n-4] == 'e' && dst[n-3] == '-' && dst[n-2] == '0' {
		dst[n-2] = dst[n-1]
		dst = dst[:n-1]
	}
	return dst, nil
}
func (e *jsonObjectEncoder) float(key string, v float64, bits int) error {
	var scratch [32]byte
	data, err := jsonFloat(v, bits, scratch[:0])
	if err != nil {
		return err
	}
	return e.add(key, data)
}
func (e *jsonObjectEncoder) AddFloat32(key string, v float32) error {
	return e.float(key, float64(v), 32)
}
func (e *jsonObjectEncoder) AddFloat64(key string, v float64) error { return e.float(key, v, 64) }
func (e *jsonObjectEncoder) AddInt(key string, v int) error         { return e.AddInt64(key, int64(v)) }
func (e *jsonObjectEncoder) AddInt8(key string, v int8) error       { return e.AddInt64(key, int64(v)) }
func (e *jsonObjectEncoder) AddInt16(key string, v int16) error     { return e.AddInt64(key, int64(v)) }
func (e *jsonObjectEncoder) AddInt32(key string, v int32) error     { return e.AddInt64(key, int64(v)) }
func (e *jsonObjectEncoder) AddUint(key string, v uint) error       { return e.AddUint64(key, uint64(v)) }
func (e *jsonObjectEncoder) AddUint8(key string, v uint8) error     { return e.AddUint64(key, uint64(v)) }
func (e *jsonObjectEncoder) AddUint16(key string, v uint16) error   { return e.AddUint64(key, uint64(v)) }
func (e *jsonObjectEncoder) AddUint32(key string, v uint32) error   { return e.AddUint64(key, uint64(v)) }
func (e *jsonObjectEncoder) AddUintptr(key string, v uintptr) error {
	return e.AddUint64(key, uint64(v))
}

type jsonArrayEncoder struct {
	buf   *bytes.Buffer
	count int
}

func (e *jsonArrayEncoder) prefix() {
	if e.count > 0 {
		e.buf.WriteString(", ")
	}
	e.count++
}
func (e *jsonArrayEncoder) append(data []byte) error { e.prefix(); e.buf.Write(data); return nil }
func (e *jsonArrayEncoder) AppendReflected(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return e.append(data)
}
func (e *jsonArrayEncoder) AppendObject(v ObjectMarshaler) error {
	start, count := e.buf.Len(), e.count
	e.prefix()
	e.buf.WriteByte('{')
	child := jsonObjectEncoder{buf: e.buf}
	if err := v.MarshalLogObject(&child); err != nil {
		e.buf.Truncate(start)
		e.count = count
		return err
	}
	e.buf.WriteByte('}')
	return nil
}
func (e *jsonArrayEncoder) AppendArray(v ArrayMarshaler) error {
	start, count := e.buf.Len(), e.count
	e.prefix()
	e.buf.WriteByte('[')
	child := jsonArrayEncoder{buf: e.buf}
	if err := v.MarshalLogArray(&child); err != nil {
		e.buf.Truncate(start)
		e.count = count
		return err
	}
	e.buf.WriteByte(']')
	return nil
}
func (e *jsonArrayEncoder) AppendBool(v bool) error {
	e.prefix()
	if v {
		e.buf.WriteString("true")
	} else {
		e.buf.WriteString("false")
	}
	return nil
}
func (e *jsonArrayEncoder) AppendString(v string) error {
	e.prefix()
	writeJSONString(e.buf, v)
	return nil
}
func (e *jsonArrayEncoder) AppendInt64(v int64) error {
	e.prefix()
	e.buf.Write(strconv.AppendInt(e.buf.AvailableBuffer(), v, 10))
	return nil
}
func (e *jsonArrayEncoder) AppendUint64(v uint64) error {
	e.prefix()
	e.buf.Write(strconv.AppendUint(e.buf.AvailableBuffer(), v, 10))
	return nil
}
func (e *jsonArrayEncoder) float(v float64, bits int) error {
	var scratch [32]byte
	data, err := jsonFloat(v, bits, scratch[:0])
	if err != nil {
		return err
	}
	return e.append(data)
}
func (e *jsonArrayEncoder) AppendFloat32(v float32) error { return e.float(float64(v), 32) }
func (e *jsonArrayEncoder) AppendFloat64(v float64) error { return e.float(v, 64) }
func (e *jsonArrayEncoder) AppendInt(v int) error         { return e.AppendInt64(int64(v)) }
func (e *jsonArrayEncoder) AppendInt8(v int8) error       { return e.AppendInt64(int64(v)) }
func (e *jsonArrayEncoder) AppendInt16(v int16) error     { return e.AppendInt64(int64(v)) }
func (e *jsonArrayEncoder) AppendInt32(v int32) error     { return e.AppendInt64(int64(v)) }
func (e *jsonArrayEncoder) AppendUint(v uint) error       { return e.AppendUint64(uint64(v)) }
func (e *jsonArrayEncoder) AppendUint8(v uint8) error     { return e.AppendUint64(uint64(v)) }
func (e *jsonArrayEncoder) AppendUint16(v uint16) error   { return e.AppendUint64(uint64(v)) }
func (e *jsonArrayEncoder) AppendUint32(v uint32) error   { return e.AppendUint64(uint64(v)) }
func (e *jsonArrayEncoder) AppendUintptr(v uintptr) error { return e.AppendUint64(uint64(v)) }

// writeJsonValues writes a complete JSON object containing custom fields.
func (w *Writer) writeJsonValues(fields []Field) error {
	w.jsonEnc = jsonObjectEncoder{buf: &w.buf}
	w.buf.WriteByte('{')
	if err := addJSONFields(&w.jsonEnc, fields); err != nil {
		return err
	}
	w.buf.WriteByte('}')
	return nil
}

func addJSONFields(enc ObjectEncoder, fields []Field) error {
	for i, field := range fields {
		if field.Key == "" {
			field.Key = fmt.Sprintf("arg%d", i)
		}
		if duration, ok := field.Any.(time.Duration); ok {
			field.Type = StringType
			field.String = duration.String()
		}
		if err := field.AddTo(enc); err != nil {
			return fmt.Errorf("field %q: %w", field.Key, err)
		}
	}
	return nil
}

func encodeFields(fields []Field) ([]byte, error) {
	var w Writer
	if err := w.writeJsonValues(fields); err != nil {
		return nil, err
	}
	return w.buf.Bytes(), nil
}

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
