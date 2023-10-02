package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Logg struct {
	*logrus.Entry
}

func GetLogger() *Logg {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.JSONFormatter{
		PrettyPrint: true,
	}

	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.DebugLevel)
	return &Logg{logrus.NewEntry(l)}
}
