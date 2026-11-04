package src

import (
	"io"
	"time"
)

func (w *ConsoleWriter) WriteLog(out io.Writer, e *Entity, caller string) ([]byte, error) {

	w.buf.Reset()
	if err := w.encodeInto(e); err != nil {
		return nil, err
	}

	w.buf.WriteString(caller)
	w.buf.WriteByte('\n')

	_, err := out.Write(w.buf.Bytes())
	return w.buf.Bytes(), err
}

func (w *ConsoleWriter) encodeInto(e *Entity) error {
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
