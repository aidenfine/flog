package src

import (
	"bytes"
	"io"
	"time"
)

// could we possibly abstract this caller arg away later?
// Would be a bit nicer to have the caller as an Entity (should compare performance though)
func WriteLog(out io.Writer, e *Entity, caller string) ([]byte, error) {
	buf, err := encode(e)
	if err != nil {
		return nil, err
	}

	// will this have perf hit?
	if caller != "" {
		buf.WriteString(caller)
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
		if message.Len() > 0 && message.Bytes()[message.Len()-1] != ' ' {
			message.WriteByte(' ')
		}
		if err := writeConsoleValues(message, e.Fields); err != nil {
			return nil, err
		}
	}
	return message, nil
}
