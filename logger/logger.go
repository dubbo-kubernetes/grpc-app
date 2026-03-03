package logger

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

var (
	// Regex to match gRPC formatting errors like %!p(...)
	formatErrorRegex = regexp.MustCompile(`%!p\([^)]+\)`)
)

// GrpcLogger filters out xDS informational logs that are incorrectly marked as ERROR
type GrpcLogger struct {
	Logger *log.Logger
}

// cleanMessage removes formatting errors from gRPC logs
// Fixes issues like: "\u003c%!p(networktype.keyType=grpc.internal.transport.networktype)\u003e": "unix"
func cleanMessage(msg string) string {
	// Replace %!p(...) patterns with a cleaner representation
	msg = formatErrorRegex.ReplaceAllStringFunc(msg, func(match string) string {
		// Extract the key from %!p(networktype.keyType=...)
		if strings.Contains(match, "networktype.keyType") {
			return `"networktype"`
		}
		// For other cases, just remove the error pattern
		return ""
	})
	// Also clean up Unicode escape sequences that appear with formatting errors
	// Replace \u003c (which is <) and \u003e (which is >) when they appear with formatting errors
	msg = strings.ReplaceAll(msg, `\u003c`, "<")
	msg = strings.ReplaceAll(msg, `\u003e`, ">")
	// Clean up patterns like <...>: "unix" to just show the value
	msg = regexp.MustCompile(`<[^>]*>:\s*"unix"`).ReplaceAllString(msg, `"networktype": "unix"`)
	return msg
}

func (l *GrpcLogger) Info(args ...interface{}) {
	msg := fmt.Sprint(args...)
	if strings.Contains(msg, "entering mode") && strings.Contains(msg, "SERVING") {
		return
	}
	msg = cleanMessage(msg)
	l.Logger.Print("INFO: ", msg)
}

func (l *GrpcLogger) Infoln(args ...interface{}) {
	msg := fmt.Sprintln(args...)
	if strings.Contains(msg, "entering mode") && strings.Contains(msg, "SERVING") {
		return
	}
	msg = cleanMessage(msg)
	l.Logger.Print("INFO: ", msg)
}

func (l *GrpcLogger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if strings.Contains(msg, "entering mode") && strings.Contains(msg, "SERVING") {
		return
	}
	msg = cleanMessage(msg)
	l.Logger.Printf("INFO: %s", msg)
}

func (l *GrpcLogger) Warning(args ...interface{}) {
	msg := cleanMessage(fmt.Sprint(args...))
	l.Logger.Print("WARNING: ", msg)
}

func (l *GrpcLogger) Warningln(args ...interface{}) {
	msg := cleanMessage(fmt.Sprintln(args...))
	l.Logger.Print("WARNING: ", msg)
}

func (l *GrpcLogger) Warningf(format string, args ...interface{}) {
	msg := cleanMessage(fmt.Sprintf(format, args...))
	l.Logger.Printf("WARNING: %s", msg)
}

func (l *GrpcLogger) Error(args ...interface{}) {
	msg := fmt.Sprint(args...)
	if strings.Contains(msg, "entering mode") && strings.Contains(msg, "SERVING") {
		return
	}
	msg = cleanMessage(msg)
	l.Logger.Print("ERROR: ", msg)
}

func (l *GrpcLogger) Errorln(args ...interface{}) {
	msg := fmt.Sprintln(args...)
	if strings.Contains(msg, "entering mode") && strings.Contains(msg, "SERVING") {
		return
	}
	msg = cleanMessage(msg)
	l.Logger.Print("ERROR: ", msg)
}

func (l *GrpcLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if strings.Contains(msg, "entering mode") && strings.Contains(msg, "SERVING") {
		return
	}
	msg = cleanMessage(msg)
	l.Logger.Printf("ERROR: %s", msg)
}

func (l *GrpcLogger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

func (l *GrpcLogger) Fatalln(args ...interface{}) {
	l.Logger.Fatal(args...)
}

func (l *GrpcLogger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

func (l *GrpcLogger) V(level int) bool {
	return level <= 0
}
