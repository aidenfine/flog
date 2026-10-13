package src

import (
	"bytes"
	"io"
	"time"
)

func WriteLog(out io.Writer, e *Entity) ([]byte, error) {
	buf, err := encode(e)
	if err != nil {
		return nil, err
	}

	buf.WriteByte('\n')
	_, err = out.Write(buf.Bytes())
	return buf.Bytes(), err
}

func encode(e *Entity) (*bytes.Buffer, error) {
	message := bytes.NewBuffer(nil)

	if e.Time != nil {
		message.WriteString(e.Time.Format(time.RFC3339))
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
		// Structured logging: serialize named fields as a JSON object.
		// fields, err := encodeFields(e.Fields)
		fields, err := encodeConsoleValues(e.Fields)
		if err != nil {
			return nil, err
		}
		if message.Len() > 0 && message.Bytes()[message.Len()-1] != ' ' {
			message.WriteByte(' ')
		}
		message.Write(fields)
	}
	return message, nil
}
