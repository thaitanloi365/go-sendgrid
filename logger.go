package sendgrid

import (
	"log"
	"os"
)

func createLogger() *logger {
	l := &logger{l: log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds)}
	return l
}

type logger struct {
	l *log.Logger
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.output("ERROR SENDGRID "+format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.output("WARN SENDGRID "+format, v...)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.output("DEBUG SENDGRID "+format, v...)
}

func (l *logger) output(format string, v ...interface{}) {
	if len(v) == 0 {
		l.l.Print(format)
		return
	}
	l.l.Printf(format, v...)
}
