package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

type LogLevel int

const (
	// Assign an incremental index with `iota`
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

type Logger struct {
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	fatalLog *log.Logger
}

func New() *Logger {
	return &Logger{
		debugLog: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lmicroseconds),
		infoLog:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds),
		warnLog:  log.New(os.Stderr, "WARN: ", log.Ldate|log.Ltime|log.Lmicroseconds),
		errorLog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds),
		fatalLog: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lmicroseconds),
	}
}

// If incase you would like to write the log to a file
func (l *Logger) SetOutput(w io.Writer) {
	l.debugLog.SetOutput(w)
	l.infoLog.SetOutput(w)
	l.warnLog.SetOutput(w)
	l.errorLog.SetOutput(w)
	l.fatalLog.SetOutput(w)
}

// The Core logging function
func (l *Logger) log(level LogLevel, v ...interface{}) {
	msg := fmt.Sprint(v...)
	switch level {
	case DEBUG:
		l.debugLog.Output(2, msg)
	case INFO:
		l.infoLog.Output(2, msg)
	case WARN:
		l.warnLog.Output(2, msg)
	case ERROR:
		l.errorLog.Output(2, msg)
	case FATAL:
		l.fatalLog.Output(2, msg)
		os.Exit(1)
	}
}

// logf - formatted log ( Just following the standard loggers )
func (l *Logger) logf(level LogLevel, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.log(level, msg)
}

func (l *Logger) Debug(v ...interface{})                 { l.log(DEBUG, v...) }
func (l *Logger) Debugf(format string, v ...interface{}) { l.logf(DEBUG, format, v...) }
func (l *Logger) Info(v ...interface{})                  { l.log(INFO, v...) }
func (l *Logger) Infof(format string, v ...interface{})  { l.logf(INFO, format, v...) }
func (l *Logger) Warn(v ...interface{})                  { l.log(WARN, v...) }
func (l *Logger) Warnf(format string, v ...interface{})  { l.logf(WARN, format, v...) }
func (l *Logger) Error(v ...interface{})                 { l.log(ERROR, v...) }
func (l *Logger) Errorf(format string, v ...interface{}) { l.logf(ERROR, format, v...) }
func (l *Logger) Fatal(v ...interface{})                 { l.log(FATAL, v...) }
func (l *Logger) Fatalf(format string, v ...interface{}) { l.logf(FATAL, format, v...) }

func (l *Logger) LogRequest(method, path, remoteAddr string, statusCode int) {
	l.Infof("Reqeust: %s %d %s %s", method, statusCode, path, remoteAddr)
}

// logs server errors with file and line info
func (l *Logger) LogServerError(err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		l.Errorf("Server Error at %s:%d: %v", file, line, err)
	} else {
		l.Errorf("Server error: %v", err)
	}
}

// logs component failure with extra context
func (l *Logger) LogFailure(component string, err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		l.Errorf("Failure in %s at %s:%d: %v", component, file, line, err)
	} else {
		l.Errorf("failure in %s: %v", component, err)
	}
}
