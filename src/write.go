package src

import (
	"bytes"
	"io"
	"time"
)

// could we possibly abstract this caller arg away later?
// Would be a bit nicer to have the caller as an Entity (should compare performance though)
func WriteLog(out io.Writer, e *Entity, caller string) ([]byte, error) {
	return WriteLogBuffer(out, new(bytes.Buffer), e, caller)
}

// WriteLogBuffer resets and reuses buf. Returned bytes are valid until buf is reused.
// Callers must synchronize access to buf and finish using it before the next call.
func WriteLogBuffer(out io.Writer, buf *bytes.Buffer, e *Entity, caller string) ([]byte, error) {
	buf.Reset()
	if err := encodeInto(buf, e); err != nil {
		return nil, err
	}

	// will this have perf hit?
	if caller != "" {
		buf.WriteString(caller)
	}

	buf.WriteByte('\n')
	_, err := out.Write(buf.Bytes())
	return buf.Bytes(), err
}

func encode(e *Entity) (*bytes.Buffer, error) {
	message := new(bytes.Buffer)
	if err := encodeInto(message, e); err != nil {
		return nil, err
	}
	return message, nil
}

func encodeInto(message *bytes.Buffer, e *Entity) error {

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
		if err := writeConsoleValues(message, e.Fields); err != nil {
			return err
		}
	}
	return nil
}
