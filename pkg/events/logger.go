package events

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
)

// Logger is the logging contract used by the watcher: Printf for regular
// messages and Error for failures. logrus and zap's SugaredLogger satisfy it
// directly; use LogAdapter for Printf-only loggers such as *log.Logger, or
// SlogLogger / StdLogger for the implementations in this file.
type Logger interface {
	Printf(format string, args ...any)
	Error(args ...any)
}

// Printfer is the minimal method set of the standard *log.Logger.
type Printfer interface {
	Printf(format string, args ...any)
}

// StdLogger is a Logger backed by the standard library log.Logger. Error
// messages are tagged with an "ERROR" prefix.
type StdLogger struct {
	*log.Logger
}

// NewStdLogger returns a StdLogger writing to w with date and time flags.
func NewStdLogger(w io.Writer) *StdLogger {
	return &StdLogger{Logger: log.New(w, "", log.LstdFlags)}
}

// Printf implements Logger.
func (l *StdLogger) Printf(format string, args ...any) {
	if l != nil && l.Logger != nil {
		l.Logger.Printf(format, args...)
	}
}

// Error implements Logger.
func (l *StdLogger) Error(args ...any) {
	if l != nil && l.Logger != nil {
		l.Logger.Output(2, "ERROR "+fmt.Sprint(args...))
	}
}

// Errorf formats an error message and forwards it to Error.
func (l *StdLogger) Errorf(format string, args ...any) {
	l.Error(fmt.Sprintf(format, args...))
}

// LogAdapter wraps a Printf-only logger (for example *log.Logger) so it can be
// used as a Logger; Error calls are routed through Printf.
type LogAdapter struct {
	Printfer Printfer
}

// Printf implements Logger.
func (l LogAdapter) Printf(format string, args ...any) {
	if l.Printfer != nil {
		l.Printfer.Printf(format, args...)
	}
}

// Error implements Logger by forwarding to Printf.
func (l LogAdapter) Error(args ...any) { l.Printf("%s", fmt.Sprint(args...)) }

// NopLogger discards every message. It is the default logger.
type NopLogger struct{}

// Printf implements Logger.
func (NopLogger) Printf(string, ...any) {}

// Error implements Logger.
func (NopLogger) Error(...any) {}

// SlogLogger adapts log/slog to the Logger interface. A zero Level logs at
// slog.LevelError, which is the level used for watcher errors.
type SlogLogger struct {
	*slog.Logger
	Level slog.Level
}

// Printf implements Logger.
func (l SlogLogger) Printf(format string, args ...any) {
	l.log(fmt.Sprintf(format, args...))
}

// Error implements Logger.
func (l SlogLogger) Error(args ...any) { l.log(fmt.Sprint(args...)) }

func (l SlogLogger) log(msg string) {
	if l.Logger == nil {
		return
	}
	level := l.Level
	if level == 0 {
		level = slog.LevelError
	}
	l.Logger.Log(context.Background(), level, msg)
}
