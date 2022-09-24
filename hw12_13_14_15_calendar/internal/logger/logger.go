package logger

import (
	"io"
	"log"
	"os"
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
	case "off":
		return LogLevel(Off)
	default:
		return LogLevel(Error)
	}
}

const (
	Debug int = iota
	Info
	Warn
	Error
	Off
)

type Logger struct {
	level                                            LogLevel
	debugPrefix, infoPrefix, warnPrefix, errorPrefix string
	logg                                             *log.Logger
}

func New(level string) *Logger {
	logg := log.New(os.Stdout, "", 0)
	return &Logger{
		level:       NewLogLevel(level),
		debugPrefix: "[DEBUG]",
		infoPrefix:  "[INFO]",
		warnPrefix:  "[WARN]",
		errorPrefix: "[ERROR]",
		logg:        logg,
	}
}

func (l *Logger) SetOutput(out io.Writer) {
	l.logg.SetOutput(out)
}

func (l Logger) Debug(msg string) {
	if l.level <= LogLevel(Debug) {
		l.logg.Printf("%s %s", l.debugPrefix, msg)
	}
}

func (l Logger) Info(msg string) {
	if l.level <= LogLevel(Info) {
		l.logg.Printf("%s %s\n", l.infoPrefix, msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.level <= LogLevel(Warn) {
		l.logg.Printf("%s %s\n", l.warnPrefix, msg)
	}
}

func (l Logger) Error(msg string) {
	if l.level <= LogLevel(Error) {
		l.logg.Printf("%s %s\n", l.errorPrefix, msg)
	}
}
