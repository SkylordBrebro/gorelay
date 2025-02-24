package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warning
	Error
	Success
)

var levelColors = map[LogLevel]string{
	Debug:   "\033[36m", // Cyan
	Info:    "\033[37m", // White
	Warning: "\033[33m", // Yellow
	Error:   "\033[31m", // Red
	Success: "\033[32m", // Green
}

const (
	colorReset = "\033[0m"
)

type Logger struct {
	file  *os.File
	debug bool
}

func New(logPath string, debug bool) (*Logger, error) {
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &Logger{
		file:  f,
		debug: debug,
	}, nil
}

func (l *Logger) Log(sender string, message string, level LogLevel) {
	timestamp := time.Now().Format("15:04:05")
	color := levelColors[level]
	logMsg := fmt.Sprintf("%s[%s | %s] %s%s", color, timestamp, sender, message, colorReset)
	plainMsg := fmt.Sprintf("[%s | %s] %s", timestamp, sender, message)

	if level == Debug && !l.debug {
		return
	}

	log.Println(logMsg)
	fmt.Fprintln(l.file, plainMsg)
}

// Convenience methods
func (l *Logger) Debug(sender string, format string, args ...interface{}) {
	l.Log(sender, fmt.Sprintf(format, args...), Debug)
}

func (l *Logger) Info(sender string, format string, args ...interface{}) {
	l.Log(sender, fmt.Sprintf(format, args...), Info)
}

func (l *Logger) Warning(sender string, format string, args ...interface{}) {
	l.Log(sender, fmt.Sprintf(format, args...), Warning)
}

func (l *Logger) Error(sender string, format string, args ...interface{}) {
	l.Log(sender, fmt.Sprintf(format, args...), Error)
}

func (l *Logger) Success(sender string, format string, args ...interface{}) {
	l.Log(sender, fmt.Sprintf(format, args...), Success)
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// SetDebug enables or disables debug logging
func (l *Logger) SetDebug(debug bool) {
	l.debug = debug
}
