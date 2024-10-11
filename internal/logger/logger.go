package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func New() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		PadLevelText:    true})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	return log
}
