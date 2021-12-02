package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type UTCFormatter struct {
	logrus.Formatter
}

func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()

	return u.Formatter.Format(e)
}

var (
	Logger        = logrus.WithFields(logrus.Fields{})
	HTTPLogFormat = fmt.Sprintf("%s\n", `{"time": ${time_unix_nano}, }`) // Annoying \n hack
)

func ConfigureLogger() {
	Logger.Logger.SetLevel(logrus.DebugLevel)
	Logger.Logger.SetFormatter(UTCFormatter{&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z",
	}})
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}
func Info(args ...interface{}) {
	Logger.Info(args...)
}
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}
func Error(args ...interface{}) {
	Logger.Error(args...)
}
