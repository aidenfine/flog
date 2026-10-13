package flog

import "github.com/aidenfine/flog/src"

// Field is a named value added to a log entry.
type Field = src.Field
type ObjectEncoder = src.ObjectEncoder
type ArrayEncoder = src.ArrayEncoder
type ObjectMarshaler = src.ObjectMarshaler
type ArrayMarshaler = src.ArrayMarshaler

// Any creates a named field, supporting custom marshalers and arbitrary values.
func Any(key string, value any) Field {
	field := src.MakeFields(value)[0]
	field.Key = key
	return field
}
