package log

import (
	"os"
	"github.com/sirupsen/logrus"
)

// Create a new instance of the logger. You can have any number of instances.
var log = logrus.New()

func Info(any interface{}) {
	log.Info(any)
}

func Warn(any interface{}) {
	log.Warn(any)
}

func Error(any interface{}) {
	log.Error(any)
}


func init() {
	log.Out = os.Stdout

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")
}
