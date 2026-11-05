package src

// map of the level to int by priority
var LogLevelToInt = map[string]int{
	"ERROR": 4,
	"WARN":  3,
	"INFO":  2,
	"DEBUG": 1,
}

func capsLevelToLower(level string) string {
	switch level {
	case "ERROR":
		return "error"
	case "DEBUG":
		return "debug"
	case "INFO":
		return "info"
	case "WARN":
		return "warn"
	}
	return ""
}
func lowerLevelToCaps(level string) string {
	switch level {
	case "error":
		return "ERROR"
	case "debug":
		return "DEBUG"
	case "info":
		return "INFO"
	case "warn":
		return "WARN"
	}
	return ""
}
