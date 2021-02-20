package logger

import (
	"io"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Output() io.Writer
	SetOutput(w io.Writer)
}

var DefaultLogger Logger

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args)
}
func Info(args ...interface{}) {
	DefaultLogger.Info(args)
}

func Warn(args ...interface{}) {
	DefaultLogger.Warn(args)
}

func Error(args ...interface{}) {
	DefaultLogger.Error(args)
}
func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args)
}
func Panic(args ...interface{}) {
	DefaultLogger.Panic(args)
}
