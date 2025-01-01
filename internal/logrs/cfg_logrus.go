package logrs

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitLogrus() {
	logFile, err := os.OpenFile("./FinMarket.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create or open log file")
	}

	logrus.SetOutput(logFile)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}
