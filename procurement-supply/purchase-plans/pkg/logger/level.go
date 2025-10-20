package logger

// Level represents the severity of a log message
type Level string

const (
	Debug Level = "DEBUG"
	Info  Level = "INFO"
	Warn  Level = "WARN"
	Error Level = "ERROR"
	Fatal Level = "FATAL"
)
