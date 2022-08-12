package logger

import (
	"fmt"
	"strings"
)

type LogLevel int

func NewLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return LogLevel(Debug)
	case "info":
		return LogLevel(Info)
	case "warn":
		return LogLevel(Warn)
	default:
		return LogLevel(Error)
	}
}

const (
	Debug int = iota
	Info
	Warn
	Error
	None
)

type Logger struct {
	level                                            LogLevel
	debugPrefix, infoPrefix, warnPrefix, errorPrefix string
}

func New(level string) *Logger {
	return &Logger{
		level:       NewLogLevel(level),
		debugPrefix: "[DEBUG]",
		infoPrefix:  "[INFO]",
		warnPrefix:  "[WARN]",
		errorPrefix: "[ERROR]",
	}
}

func (l Logger) Debug(msg string) {
	if l.level >= LogLevel(Debug) {
		fmt.Printf("%s %s", l.debugPrefix, msg)
	}
}

func (l Logger) Info(msg string) {
	if l.level >= LogLevel(Info) {
		fmt.Printf("%s %s\n", l.infoPrefix, msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.level >= LogLevel(Warn) {
		fmt.Printf("%s %s\n", l.warnPrefix, msg)
	}
}

func (l Logger) Error(msg string) {
	if l.level != LogLevel(None) {
		fmt.Printf("%s %s\n", l.errorPrefix, msg)
	}
}
