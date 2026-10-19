package flog

import (
	"fmt"
	"strings"
)

// formatMessage substitutes positional values for {}. Explicit fields remain
// fields; unused values retain the existing appended-value behavior.
func formatMessage(msg string, args []any) (string, []any) {
	if !strings.ContainsAny(msg, "{}") {
		return msg, args
	}
	var out strings.Builder
	out.Grow(len(msg))
	var remaining []any
	next := 0
	for i := 0; i < len(msg); {
		if i+1 < len(msg) {
			switch msg[i : i+2] {
			case "{{", "}}":
				out.WriteByte(msg[i])
				i += 2
				continue
			case "{}":
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
					out.WriteString("{}")
				}
				i += 2
				continue
			}
		}
		out.WriteByte(msg[i])
		i++
	}
	if len(remaining) == 0 {
		return out.String(), args[next:]
	}
	return out.String(), append(remaining, args[next:]...)
}
