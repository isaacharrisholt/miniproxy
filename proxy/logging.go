package proxy

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type logLine struct {
	line        string
	serviceName string
	level       string
}

func (l logLine) String() string {
	return fmt.Sprintf("[%s] %s: %s", strings.ToUpper(l.level)+strings.Repeat(" ", 5-len(l.level)), l.serviceName, l.line)
}

func stdOutReceiver() {
	for {
		select {
		case line := <-stdOutChan:
			log.Println(line)
		}
	}
}

type Logger interface {
	debug(message string)
	info(message string)
	error(message string)
}

type channelLogger struct {
	channel     chan<- logLine
	serviceName string
	level string
}

func (l *channelLogger) debug(message string) {
	if l.level != "debug" {
		return
	}
	l.channel <- logLine{message, l.serviceName, "debug"}
}

func (l *channelLogger) info(message string) {
	if l.level == "error" {
		return
	}
	l.channel <- logLine{message, l.serviceName, "info"}
}

func (l *channelLogger) error(message string) {
	l.channel <- logLine{message, l.serviceName, "error"}
}

func NewChannelLogger(serviceName string) *channelLogger {
	level := os.Getenv("LOG_LEVEL")
	if level != "info" && level != "debug" && level != "error" {
		level = "info"
	}
	return &channelLogger{
		channel: stdOutChan,
		level: level,
		serviceName: serviceName,
	}
}