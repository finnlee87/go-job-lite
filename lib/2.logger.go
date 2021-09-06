package lib

import (
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// Logger .
var Logger *logrus.Logger

// NewLogger is Constructed
func init() {
	Logger = logrus.New()
	Logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	Logger.SetReportCaller(true)
	mode := Config.Get("mode")
	if mode == "alpha" || mode == "beta" || mode == "staging" || mode == "release" {
		if mode == "alpha" || mode == "beta" {
			Logger.SetLevel(logrus.DebugLevel)
		} else {
			Logger.SetLevel(logrus.InfoLevel)
		}
		logPath := Config.GetString("logPath")
		logf, err := rotatelogs.New(
			logPath+"access.log.%Y%m%d",
			rotatelogs.WithMaxAge(24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			fmt.Errorf("failed to create rotatelogs: %s", err.Error())
			return
		}
		Logger.Out = logf
	} else {
		Logger.SetLevel(logrus.DebugLevel)
		Logger.Out = os.Stdout
	}
}
