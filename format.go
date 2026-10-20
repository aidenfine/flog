package flog

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// formatMessage substitutes positional values for separator, which must contain
// exactly two UTF-8 characters. For distinct characters, doubling either escapes it. Explicit
// fields remain fields; unused values retain the existing appended-value behavior.
// It panics if separator is invalid.
func formatMessage(msg string, args []any, separator string) (string, []any) {
	if !utf8.ValidString(separator) || utf8.RuneCountInString(separator) != 2 {
		panic("flog: separator must contain exactly two characters")
	}
	_, firstSize := utf8.DecodeRuneInString(separator)
	first, second := separator[:firstSize], separator[firstSize:]
	if !strings.ContainsAny(msg, separator) {
		return msg, args
	}
	var out strings.Builder
	out.Grow(len(msg))
	var remaining []any
	next := 0
	for i := 0; i < len(msg); {
		switch {
		case first != second && strings.HasPrefix(msg[i:], first+first):
			out.WriteString(first)
			i += 2 * len(first)
			continue
		case first != second && strings.HasPrefix(msg[i:], second+second):
			out.WriteString(second)
			i += 2 * len(second)
			continue
		case strings.HasPrefix(msg[i:], separator):
			for next < len(args) {
				if _, ok := args[next].(Field); !ok {
					break
				}
				remaining = append(remaining, args[next])
				next++
			}
			if next < len(args) {
				fmt.Fprint(&out, args[next])
				next++
			} else {
				out.WriteString(separator)
			}
			i += len(separator)
			continue
		}
		out.WriteByte(msg[i])
		i++
	}
	if len(remaining) == 0 {
		return out.String(), args[next:]
	}
	return out.String(), append(remaining, args[next:]...)
}
