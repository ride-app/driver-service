package logger

import (
	"os"

	"github.com/ride-app/driver-service/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	// Trace(args ...interface{})
	// Tracef(format string, args ...interface{})
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
	logger *zap.SugaredLogger
}

func New() *LogrusLogger {
	encoderConfig := zap.NewProductionEncoderConfig()

	if !config.Env.Production {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.LevelKey = "severity"
	encoderConfig.MessageKey = "message"
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	if !config.Env.Production {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       !config.Env.Production,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderConfig,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	if !config.Env.Production {
		zapConfig.Encoding = "console"
	}

	if config.Env.Debug {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	z := zap.Must(zapConfig.Build(zap.AddCallerSkip(1))).Sugar()

	return &LogrusLogger{
		logger: z,
	}
}

// func (l *LogrusLogger) Trace(args ...interface{}) {
// 	l.logger.Trace(args...)
// }

// func (l *LogrusLogger) Tracef(format string, args ...interface{}) {
// 	l.logger.Tracef(format, args...)
// }

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
		logger: l.logger.With(key, value),
	}
}

func (l *LogrusLogger) WithFields(fields map[string]string) Logger {
	return &LogrusLogger{
		logger: l.logger.With(fields),
	}
}

func (l *LogrusLogger) WithError(err error) Logger {
	return &LogrusLogger{
		logger: l.logger.With(zap.Error(err)),
	}
}
