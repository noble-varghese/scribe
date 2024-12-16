// internal/logger/logger.go
package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func init() {
	// Custom log format with timestamp, level, and message
	flags := log.Lmsgprefix

	infoLogger = log.New(os.Stdout, getPrefix("INFO"), flags)
	warningLogger = log.New(os.Stdout, getPrefix("WARN"), flags)
	errorLogger = log.New(os.Stderr, getPrefix("ERROR"), flags)
}

func getPrefix(level string) string {
	return fmt.Sprintf("%-5s | ", level)
}

func formatMessage(v ...interface{}) string {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	return fmt.Sprintf("%s | %s", timestamp, fmt.Sprint(v...))
}

func Info(v ...interface{}) {
	infoLogger.Println(formatMessage(v...))
}

func Infof(format string, v ...interface{}) {
	infoLogger.Println(formatMessage(fmt.Sprintf(format, v...)))
}

func Warn(v ...interface{}) {
	warningLogger.Println(formatMessage(v...))
}

func Warnf(format string, v ...interface{}) {
	warningLogger.Println(formatMessage(fmt.Sprintf(format, v...)))
}

func Error(v ...interface{}) {
	errorLogger.Println(formatMessage(v...))
}

func Errorf(format string, v ...interface{}) {
	errorLogger.Println(formatMessage(fmt.Sprintf(format, v...)))
}

func Fatal(v ...interface{}) {
	errorLogger.Fatal(formatMessage(v...))
}

func Fatalf(format string, v ...interface{}) {
	errorLogger.Fatal(formatMessage(fmt.Sprintf(format, v...)))
}
