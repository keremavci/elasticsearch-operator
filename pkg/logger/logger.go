package logger

import (
	"encoding/json"
	"os"

	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
)

var log logr.Logger

//Debugf logs messages at level 2
func Debugf(format string, objects ...interface{}) {
	logrus.Debugf(format, objects...)
}

//IsDebugEnabled returns true if loglevel is 2
func IsDebugEnabled() bool {
	return logrus.GetLevel() == logrus.DebugLevel
}

func Warnf(format string, objects ...interface{}) {
	logrus.Warnf(format, objects...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Infof(format string, objects ...interface{}) {
	logrus.Infof(format, objects...)
}

func init() {
	level := os.Getenv("LOG_LEVEL")
	parsed, err := logrus.ParseLevel(level)
	if err != nil {
		parsed = logrus.InfoLevel
		logrus.Warnf("Unable to parse loglevel %q", level)
	}
	logrus.SetLevel(parsed)
}

//DebugObject pretty prints the given object
func DebugObject(sprintfMessage string, object interface{}) {
	if IsDebugEnabled() && object != nil {
		pretty, err := json.MarshalIndent(object, "", "  ")
		if err != nil {
			logrus.Debugf("Error marshalling object %v for debug log: %v", object, err)
		}
		logrus.Debugf(sprintfMessage, string(pretty))
	}
}
