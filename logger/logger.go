package logger

import (
	"time"

	"github.com/ride-app/driver-service/config"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]string) Logger
	WithError(err error) Logger
}

type LogrusLogger struct {
	logger *logrus.Entry
}

func New() *LogrusLogger {
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	})

	logrus.SetLevel(logrus.InfoLevel)

	if config.Env.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return &LogrusLogger{
		logger: logrus.WithFields(logrus.Fields{
			"PROJECT_ID": config.Env.ProjectID,
		}),
	}
}

// func (l *LogrusLogger) WithError(err error) *LogrusLogger {
// 	return l.logger.WithField("error", err.Error())
// }

func (l *LogrusLogger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *LogrusLogger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *LogrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *LogrusLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *LogrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	return &LogrusLogger{
		logger: l.logger.WithField(key, value),
	}
}

func (l *LogrusLogger) WithFields(fields map[string]string) Logger {
	logFields := make(logrus.Fields)
	for k, v := range fields {
		logFields[k] = v
	}
	return &LogrusLogger{
		logger: l.logger.WithFields(logFields),
	}
}

func (l *LogrusLogger) WithError(err error) Logger {
	return &LogrusLogger{
		logger: l.logger.WithError(err),
	}
}
