package g

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger(level string) {
	switch level {
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	}

	Logger.SetOutput(os.Stdout)
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15-04-05",
	})
}
