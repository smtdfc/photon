package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	target string
	file   *os.File
}

func New(target string) *Logger {
	l := &Logger{target: target}

	if target != "stdout" {
		f, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Can't open log file %s: %v", target, err)
		}
		l.file = f
	}

	return l
}

func (l *Logger) Info(msg string) {
	l.write("INFO", msg)
}

func (l *Logger) Success(msg string) {
	l.write("SUCCESS", msg)
}

func (l *Logger) Warn(msg string) {
	l.write("WARN", msg)
}

func (l *Logger) Error(msg string) {
	l.write("ERROR", msg)
}

func (l *Logger) write(level, msg string) {
	finalMsg := fmt.Sprintf("[%s] %s\n", level, msg)
	if l.target == "stdout" {
		fmt.Print(finalMsg)
	} else {
		if l.file != nil {
			l.file.WriteString(finalMsg)
		}
	}
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
