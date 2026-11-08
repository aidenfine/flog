package src

import (
	"fmt"
	"io"
	"strings"
	"time"
)

func (w *Writer) WriteLog(out io.Writer, e *Entity, caller, encoding string) ([]byte, error) {

	w.buf.Reset()
	switch encoding {
	case "", "console":
		if err := w.consoleEncode(e); err != nil {
			return nil, err
		}
		w.buf.WriteString(caller)
	case "json":
		if err := w.jsonEncode(e, caller); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported encoding %q", encoding)
	}

	w.buf.WriteByte('\n')

	_, err := out.Write(w.buf.Bytes())
	return w.buf.Bytes(), err
}

func (w *Writer) consoleEncode(e *Entity) error {
	message := &w.buf

	if e.Time != nil {
		layout := e.TimeFormat
		if layout == "" {
			layout = time.RFC3339
		}
		message.Grow(64)
		message.Write(e.Time.AppendFormat(message.AvailableBuffer(), layout))
		message.WriteByte(' ')
	}

	if e.LogLevel != "" {
		message.WriteByte('[')
		message.WriteString(e.LogLevel)
		message.WriteString("] ")
	}

	if e.LoggerName != "" {
		message.WriteByte('[')
		message.WriteString(e.LoggerName)
		message.WriteString("] - ")
	}
	if e.Msg != "" {
		message.WriteString(e.Msg)
	}

	if len(e.Fields) > 0 {
		if message.Len() > 0 && message.Bytes()[message.Len()-1] != ' ' {
			message.WriteByte(' ')
		}
		if err := w.writeConsoleValues(e.Fields); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) jsonEncode(e *Entity, caller string) error {
	w.jsonEnc = jsonObjectEncoder{buf: &w.buf}
	w.buf.WriteByte('{')
	metadata := []struct{ key, value string }{
		{"msg", e.Msg},
		{"level", e.LogLevel},
		{"logger", e.LoggerName},
		{"caller", strings.TrimPrefix(caller, " - ")},
	}
	for _, field := range metadata {
		if field.value != "" || field.key == "msg" {
			if err := w.jsonEnc.AddString(field.key, field.value); err != nil {
				return err
			}
		}
	}
	if e.Time != nil {
		// layout := e.TimeFormat
		// if layout == "" {
		// layout = time.RFC3339
		// }
		if err := w.jsonEnc.AddString("time", RFC3339TimeEncoder(*e.Time)); err != nil {
			return err
		}
	}
	if len(e.Fields) > 0 {
		if err := w.jsonEnc.AddObject("fields", jsonFields(e.Fields)); err != nil {
			return err
		}
	}
	w.buf.WriteByte('}')
	return nil
}

type jsonFields []Field

func (fields jsonFields) MarshalLogObject(enc ObjectEncoder) error {
	return addJSONFields(enc, fields)
}
