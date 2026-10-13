package src

import (
	"time"
)

// TODO: if nothing more will be done with entity in this file, we should remove it.
type Entity struct {
	Msg        string  // message
	Fields     []Field // list of fields
	LogLevel   string
	Time       *time.Time
	LoggerName string
}
