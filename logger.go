package photon

import (
	"log"
	"os"
)

type Logger struct {
	Target string
}

func (l *Logger) Info(msg string) {
	log.Printf("[INFO ][@modules/%s] %s\n", l.Target, msg)
}

func (l *Logger) Warn(msg string) {
	log.Printf("[WARN ][@modules/%s] %s\n", l.Target, msg)
}

func (l *Logger) Error(msg string) {
	log.Printf("[ERROR][@modules/%s] %s\n", l.Target, msg)
}

func (l *Logger) Debug(msg string) {
	log.Printf("[DEBUG][@modules/%s] %s\n", l.Target, msg)
}

func (l *Logger) Fatal(msg string) {
	log.Printf("[FATAL][@modules/%s] %s\n", l.Target, msg)
	os.Exit(1)
}


func CreateLogger(target string) *Logger {
	return &Logger{
		Target: target,
	}
}